package models

type Config struct {
	FilePath    string
	ApiEndpoint string
	Separator   string
	BufferSize  int
	Rps         int
	ChannelSize int
	DebugMode   bool
	Mapping     Mapping
}

func (c *Config) Initialize() {
	c.ChannelSize = c.Rps * 2
}

func (c *Config) Validate() bool {
	return c.FilePath != "" && c.ApiEndpoint != "" && c.Mapping.Validate()
}
