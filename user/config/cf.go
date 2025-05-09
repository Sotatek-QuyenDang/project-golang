package config

import (
	"os"
)

type Config struct {
	USER     string
	PASSWORD string
	NAME     string
	HOST     string
	PORT     string
	SSLMODE  string
}

func LoadConfig() Config {
	return Config{
		USER:     os.Getenv("DB_USER"),
		PASSWORD: os.Getenv("DB_PASSWORD"),
		NAME:     os.Getenv("DB_NAME"),
		HOST:     os.Getenv("DB_HOST"),
		PORT:     os.Getenv("DB_PORT"),
		SSLMODE:  os.Getenv("DB_SSLMODE"),
	}
}
