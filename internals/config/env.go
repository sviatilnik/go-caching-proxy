package config

import (
	"os"
)

type Config struct {
	Port   string
	Target string
}

func NewConfig() *Config {
	conf := &Config{}

	conf.init()

	return conf
}

func (c *Config) init() {
	c.Target = os.Getenv("TARGET")
	c.Port = os.Getenv("PORT")
}
