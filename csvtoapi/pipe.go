package csvtoapi

import (
	"context"
	"cvloader/models"
	"golang.org/x/time/rate"
	"log"
	"net/url"
	"sync"
)

type pipe struct {
	config          models.Config
	wg              *sync.WaitGroup
	rawChan         chan models.Raw
	queryParamsChan chan QueryParams
	urlChan         chan url.URL
}

func NewPipe(config models.Config, wg *sync.WaitGroup) *pipe {
	return &pipe{
		config:          config,
		wg:              wg,
		queryParamsChan: make(chan QueryParams, config.ChannelSize),
		urlChan:         make(chan url.URL, config.ChannelSize),
	}
}

func (p *pipe) Run() chan error {
	return p.consume().prepareLine().makeUrl().send()
}

func (p *pipe) consume() *pipe {
	p.rawChan = consume(p.wg, p.config)
	return p
}

func (p *pipe) prepareLine() *pipe {
	go func() {
		p.wg.Add(1)
		defer p.wg.Done()
		defer close(p.queryParamsChan)

		wg := sync.WaitGroup{}
		for raw := range p.rawChan {
			wg.Add(1)
			go func() {
				defer wg.Done()
				line, err := prepareLine(p.config, raw)
				if err != nil {
					log.Println(err)
				}
				p.queryParamsChan <- line
			}()
		}
		wg.Wait()
	}()
	return p
}

func (p *pipe) makeUrl() *pipe {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(p.urlChan)

		ctx := context.Background()
		limiter := rate.NewLimiter(rate.Limit(p.config.Rps), 1)

		wg := sync.WaitGroup{}
		for queryParams := range p.queryParamsChan {
			_ = limiter.Wait(ctx)

			wg.Add(1)
			go func() {
				defer wg.Done()
				targetUrl, err := makeUrl(p.config, queryParams)
				if err != nil {
					log.Println(err)
				}
				p.urlChan <- *targetUrl
			}()
		}
		wg.Wait()
	}()

	return p
}

func (p *pipe) send() chan error {
	errChan := make(chan error)
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		defer close(errChan)

		wg := sync.WaitGroup{}
		for targetUrl := range p.urlChan {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := send(p.config, targetUrl)
				if err != nil {
					errChan <- err
				}
			}()
		}
		wg.Wait()
	}()

	return errChan
}
