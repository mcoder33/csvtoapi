package main

import (
	"fmt"
	"strings"
)

const (
	mapElemSeparator     = ","
	mapKeyValueSeparator = ":"
)

type Mapping struct {
	colMap  string
	storage map[string]string
}

func (m *Mapping) queryParam(name string) (string, error) {
	elem, ok := m.storage[name]
	if !ok {
		return "", fmt.Errorf("no element named %q storage %v", name, m.storage)
	}
	return elem, nil
}

func (m *Mapping) validate() bool {
	return m.colMap != "" && strings.Contains(m.colMap, mapKeyValueSeparator)
}

func (m *Mapping) parse() error {
	elemsArr := strings.Split(m.colMap, mapElemSeparator)
	m.storage = make(map[string]string, len(elemsArr))
	if len(elemsArr) == 0 {
		return fmt.Errorf("storage source shoul be grater than 0: %s", m.colMap)
	}

	for _, keyValueString := range elemsArr {
		arr := strings.Split(keyValueString, mapKeyValueSeparator)
		if len(arr) != 2 {
			return fmt.Errorf("invalid format for storage: %s %v", keyValueString, m)
		}
		m.storage[arr[0]] = arr[1]
	}

	return nil
}
