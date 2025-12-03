package models

import (
	"fmt"
	"log"
	"strings"
)

const (
	mapElemSeparator     = ","
	mapKeyValueSeparator = ":"
)

type Mapping struct {
	ColMap  string
	Storage map[string]string
}

func (m *Mapping) QueryParam(name string) string {
	elem, ok := m.Storage[name]
	if !ok {
		return ""
	}
	return elem
}

func (m *Mapping) Validate() bool {
	err := m.Parse()
	if err != nil {
		log.Printf("Error parsing mapping: %v", err)
	}
	return err == nil && m.ColMap != "" && strings.Contains(m.ColMap, mapKeyValueSeparator)
}

func (m *Mapping) Parse() error {
	elemsArr := strings.Split(m.ColMap, mapElemSeparator)
	m.Storage = make(map[string]string, len(elemsArr))
	if len(elemsArr) == 0 {
		return fmt.Errorf("Storage source shoul be grater than 0: %s", m.ColMap)
	}

	for _, keyValueString := range elemsArr {
		arr := strings.Split(keyValueString, mapKeyValueSeparator)
		if len(arr) != 2 {
			return fmt.Errorf("invalid format for Storage: %s %v", keyValueString, m)
		}
		m.Storage[arr[0]] = arr[1]
	}

	return nil
}
