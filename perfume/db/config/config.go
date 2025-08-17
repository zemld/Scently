package config

import (
	"os"
)

type configKeyType struct{}

var ConfigKey = configKeyType{}

type Config struct {
	User     string
	Db       string
	Host     string
	Password string
}

func NewConfig() *Config {
	return &Config{User: os.Getenv("POSTGRES_USER"), Db: os.Getenv("POSTGRES_DB"), Host: os.Getenv("POSTGRES_HOST"), Password: getPassword()}
}

func (c *Config) GetConnectionString() string {
	return "postgres//:" + c.User + ":" + c.Password + "@" + c.Host + ":5432/" + c.Db
}

func getPassword() string {
	return ""
}
