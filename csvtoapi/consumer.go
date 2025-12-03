package csvtoapi

import (
	"bufio"
	"cvloader/models"
	"cvloader/utils"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

func consume(wg *sync.WaitGroup, config models.Config) chan models.Raw {
	file, err := os.Open(config.FilePath)
	if err != nil {
		log.Fatal(err)
	}

	rawChan := make(chan models.Raw, config.ChannelSize)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(rawChan)
		defer file.Close()

		r := bufio.NewReaderSize(file, config.BufferSize)
		raw := models.Raw{}

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

			line = utils.CleanString(line, `"`, `'`, "\n", "\r")
			if raw.Headers == nil {
				raw.Headers = strings.Split(line, config.Separator)
				continue
			}
			raw.Elems = strings.Split(line, config.Separator)

			if len(raw.Headers) != len(raw.Elems) {
				log.Printf("\n[Error] LINE WAS IGNORED!!! %s \n!!! [>>>] Check the separator!!! header have %d elements; row have %d elements;\n", line, len(raw.Headers), len(raw.Elems))
				continue
			}

			rawChan <- raw
		}
	}()

	return rawChan
}
