package main

import (
	"log"
	"net/url"
	"sync"
)

type RunnerService struct {
	config          Config
	wg              *sync.WaitGroup
	rawChan         chan Raw
	queryParamsChan chan QueryParams
	urlChan         chan url.URL
}

func NewRunnerService(rawChan chan Raw, config Config, wg *sync.WaitGroup) *RunnerService {
	return &RunnerService{
		config:          config,
		wg:              wg,
		rawChan:         rawChan,
		queryParamsChan: make(chan QueryParams),
		urlChan:         make(chan url.URL),
	}
}

func (rs *RunnerService) PrepareLine() {
	go func() {
		rs.wg.Add(1)
		defer rs.wg.Done()
		defer close(rs.queryParamsChan)

		wg := sync.WaitGroup{}
		for raw := range rs.rawChan {
			wg.Add(1)
			go func() {
				defer wg.Done()
				line, err := PrepareLine(rs.config, raw)
				if err != nil {
					log.Println(err)
				}
				rs.queryParamsChan <- line
			}()
		}
		wg.Wait()
	}()
}

func (rs *RunnerService) MakeUrl() {
	go func() {
		rs.wg.Add(1)
		defer rs.wg.Done()
		defer close(rs.urlChan)

		wg := sync.WaitGroup{}
		for queryParams := range rs.queryParamsChan {
			wg.Add(1)
			go func() {
				defer wg.Done()
				targetUrl, err := MakeUrl(rs.config, queryParams)
				if err != nil {
					log.Println(err)
				}
				rs.urlChan <- *targetUrl
			}()
		}
		wg.Wait()
	}()
}

func (rs *RunnerService) Send() {
	go func() {
		rs.wg.Add(1)
		defer rs.wg.Done()

		wg := sync.WaitGroup{}
		for targetUrl := range rs.urlChan {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := Send(targetUrl)
				if err != nil {
					log.Println(err)
				}
			}()
		}
		wg.Wait()
	}()
}
