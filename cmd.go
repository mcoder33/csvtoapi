package main

import (
	"bufio"
	"cvloader/csvtoapi"
	"cvloader/models"
	"flag"
	"io"
	"log"
	"os"
	"strings"
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

	if !config.Validate() {
		flag.PrintDefaults()
		return
	}

	err := config.Mapping.Parse()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(config.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := bufio.NewReaderSize(file, config.BufferSize)
	raw := models.Raw{}

	wg := &sync.WaitGroup{}
	rawChan := make(chan models.Raw, config.ChannelSize)

	errChan := csvtoapi.NewPipe(rawChan, config, wg).Run()
	go func() {
		for err := range errChan {
			log.Println(err)
		}
	}()

	canContinue := true
	for canContinue {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				canContinue = false
			}
			if len(line) == 0 {
				break
			}
		}

		line = CleanString(line, `"`, `'`, "\n", "\r")
		if raw.Headers == nil {
			raw.Headers = strings.Split(line, config.Separator)
			continue
		}
		raw.Elems = strings.Split(line, config.Separator)

		rawChan <- raw
	}
	close(rawChan)

	wg.Wait()
}
