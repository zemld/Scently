package config

import (
	"os"
)

type Config struct {
	ConnectionString string
}

func NewConfig() *Config {
	return &Config{ConnectionString: os.Getenv("DB_URL")}
}
