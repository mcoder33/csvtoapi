package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

type QueryParams map[string]string

func PrepareLine(config Config, raw Raw) (QueryParams, error) {
	result := make(QueryParams, len(raw.elems))
	for i, elem := range raw.elems {
		colNameSource := raw.headers[i]
		if colNameSource == "" {
			return nil, fmt.Errorf("invalid queryParam %s elem %s", raw.headers[i], elem)
		}
		name, err := config.mapping.queryParam(colNameSource)
		if err != nil {
			return nil, fmt.Errorf("invalid parse line %s; colNameSource: %s; error %w", strings.Join(raw.elems, config.separator), colNameSource, err)
		}
		result[name] = elem
	}
	return result, nil
}

func MakeUrl(config Config, qp QueryParams) (*url.URL, error) {
	targetUrl, err := url.Parse(config.apiEndpoint)
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

func Send(url url.URL) error {
	log.Println(url.String())
	//resp, err := http.Get(url.String())
	//if err != nil {
	//	return fmt.Errorf("http get error: %w", err)
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	body, _ := io.ReadAll(resp.Body)
	//	return fmt.Errorf("http status: %s body: %s", resp.Status, body)
	//}

	return nil
}
