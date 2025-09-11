package csvtoapi

import (
	"cvloader/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type QueryParams map[string]string

func prepareLine(config models.Config, raw models.Raw) (QueryParams, error) {
	result := make(QueryParams, len(raw.Elems))
	for i, elem := range raw.Elems {
		colNameSource := raw.Headers[i]
		if colNameSource == "" {
			return nil, fmt.Errorf("invalid queryParam %s elem %s", raw.Headers[i], elem)
		}
		name, err := config.Mapping.QueryParam(colNameSource)
		if err != nil {
			return nil, fmt.Errorf("invalid parse line %s; colNameSource: %s; error %w", strings.Join(raw.Elems, config.Separator), colNameSource, err)
		}
		result[name] = elem
	}
	return result, nil
}

func makeUrl(config models.Config, qp QueryParams) (*url.URL, error) {
	targetUrl, err := url.Parse(config.ApiEndpoint)
	if err != nil {
		return nil, err
	}
	q := targetUrl.Query()
	for k, v := range qp {
		q.Set(k, v)
	}
	targetUrl.RawQuery = q.Encode()
	return targetUrl, nil
}

func send(config models.Config, url url.URL) error {
	if config.DebugMode {
		log.Printf("[DEBUG MODE] Sending URL: %s\n", url.String())
		return nil
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return fmt.Errorf("http get error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: url: %s\nstatus: %s\nbody: %s\n", url.String(), resp.Status, body)
	}
	log.Printf("success:\n url: %s\nstatus: %s\nbody: %s\n", url.String(), resp.Status, body)

	return nil
}
