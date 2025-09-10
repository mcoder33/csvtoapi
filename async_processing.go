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

		for raw := range rs.rawChan {
			line, err := PrepareLine(rs.config, raw)
			if err != nil {
				log.Println(err)
			}
			rs.queryParamsChan <- line
		}
	}()
}

func (rs *RunnerService) MakeUrl() {
	go func() {
		rs.wg.Add(1)
		defer rs.wg.Done()
		defer close(rs.urlChan)

		for queryParams := range rs.queryParamsChan {
			targetUrl, err := MakeUrl(rs.config, queryParams)
			if err != nil {
				log.Println(err)
			}
			rs.urlChan <- *targetUrl
		}
	}()
}

func (rs *RunnerService) Send() {
	go func() {
		rs.wg.Add(1)
		defer rs.wg.Done()

		for targetUrl := range rs.urlChan {
			err := Send(targetUrl)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}
