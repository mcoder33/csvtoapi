package main

type Config struct {
	filePath    string
	apiEndpoint string
	separator   string
	bufferSize  int
	rps         int
	mapping     Mapping
}

func (c *Config) validate() bool {
	return c.filePath != "" && c.apiEndpoint != "" && c.mapping.validate()
}
