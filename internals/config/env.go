package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrEmptyTarget = errors.New("Empty target in config")
var ErrEmptyPattern = errors.New("Empty pattern in config")
var ErrInvalidCacheTTL = errors.New("Invalid cache TTL in config")

type Config struct {
	Addr    string
	Target  string
	Pattern string
	TTL     int
}

func NewConfig() (*Config, error) {
	conf := &Config{}

	err := conf.init()
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) init() error {
	c.Addr = c.getEnvOrDefault("ADDRESS", "localhost:8080")

	c.Target = c.getEnvOrDefault("TARGET", "")
	if c.Target == "" {
		return ErrEmptyTarget
	}

	c.Pattern = c.getEnvOrDefault("PATTERN", "")
	if c.Pattern == "" {
		return ErrEmptyPattern
	}

	_, err := regexp.Compile(c.Pattern)
	if err != nil {
		return fmt.Errorf("invalid proxy pattern in config: %w", err)
	}

	ttl := c.getEnvOrDefault("CACHE_TTL", "3600")

	convertedTTL, err := strconv.Atoi(ttl)
	if err != nil {
		return ErrInvalidCacheTTL
	}

	c.TTL = convertedTTL

	return nil
}

func (c *Config) getEnvOrDefault(envKey string, def string) string {
	envVal := strings.TrimSpace(os.Getenv(envKey))

	if envVal == "" {
		return def
	}

	return envVal
}
