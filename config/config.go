package config

import (
	"os"
)

type Config struct {
	CDNHost string
}

func LoadConfig() *Config {
	cdnHost := os.Getenv("CDN_HOST")
	if cdnHost == "" {
		cdnHost = "cdn.exapmple.com"
	}
	return &Config{
		CDNHost: cdnHost,
	}
}
