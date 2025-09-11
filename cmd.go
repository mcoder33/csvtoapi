package main

import (
	"cvloader/csvtoapi"
	"cvloader/models"
	"flag"
	"log"
	"sync"
)

const (
	flagApiEndpointHelp = "Target API endpoint"
	flagFileNameHelp    = "Put here the full file path"
	flagMappingHelp     = "Field colMap like key:value,key2:value2"
	flagBufferSizeHelp  = "Buffer size for line"
	flagSeparatorHelp   = "Separator for words in line"
	flagRpsHelp         = "RPS - good speed for you"
	flagDebugModeHelp   = "Debug mode instead sending real requests just log"
)

func main() {
	config := models.Config{}

	flag.StringVar(&config.ApiEndpoint, "api", "", flagApiEndpointHelp)
	flag.StringVar(&config.Mapping.ColMap, "colMap", "", flagMappingHelp)
	flag.StringVar(&config.FilePath, "file", "", flagFileNameHelp)
	flag.StringVar(&config.Separator, "sep", ",", flagSeparatorHelp)
	flag.IntVar(&config.BufferSize, "bs", 1<<10, flagBufferSizeHelp)
	flag.IntVar(&config.Rps, "rps", 10, flagRpsHelp)
	flag.BoolVar(&config.DebugMode, "debugMode", false, flagDebugModeHelp)
	flag.Parse()

	config.Initialize()

	//TODO: переделать на возврат ошибки вместо bool
	if !config.Validate() {
		flag.PrintDefaults()
		return
	}

	err := config.Mapping.Parse()
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	//TODO: пронинуть везде контекст, добавить таймауты для хттп клиента
	//TODO: рассмотреть политику остановки при первой ошибке
	errChan := csvtoapi.NewPipe(config, wg).Run()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range errChan {
			log.Println(err)
		}
	}()

	wg.Wait()
}
