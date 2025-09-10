package main

import (
	"bufio"
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
	flagMappingHelp     = "Field raw like key:value,key2:value2"
	flagBufferSizeHelp  = "Buffer size for line"
	flagSeparatorHelp   = "Separator for words in line"
	flagRpsHelp         = "RPS - good speed for you"
)

func main() {
	config := Config{}

	flag.StringVar(&config.apiEndpoint, "api", "", flagApiEndpointHelp)
	flag.StringVar(&config.mapping.raw, "raw", "", flagMappingHelp)
	flag.StringVar(&config.filePath, "file", "", flagFileNameHelp)
	flag.StringVar(&config.separator, "sep", ",", flagSeparatorHelp)
	flag.IntVar(&config.bufferSize, "bs", 1<<10, flagBufferSizeHelp)
	flag.IntVar(&config.rps, "rps", 10, flagRpsHelp)
	flag.Parse()

	if !config.validate() {
		flag.PrintDefaults()
		return
	}

	err := config.mapping.parse()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(config.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := bufio.NewReaderSize(file, config.bufferSize)
	raw := Raw{}

	wg := &sync.WaitGroup{}
	rawChan := make(chan Raw, config.rps*2)

	rs := NewRunnerService(rawChan, config, wg)
	rs.PrepareLine()
	rs.MakeUrl()
	rs.Send()

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
		if raw.headers == nil {
			cs := CleanString(line, rune(config.separator[0]))
			raw.headers = strings.Split(cs, config.separator)
		}
		raw.elems = strings.Split(line, config.separator)

		rawChan <- raw
	}
	close(rawChan)

	wg.Wait()
}
