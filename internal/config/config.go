package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Server struct {
		ReadHeaderTimeoutSeconds time.Duration
		ReadTimeoutSeconds       time.Duration
		WriteTimeoutSeconds      time.Duration
	}
	PackSizes []int
}

// New / Init / Get
func New(path string) (*Config, error) {
	conf := &Config{}
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err := json.Unmarshal(b, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
