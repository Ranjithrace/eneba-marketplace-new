package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func FromEnvironment() (Config, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{DatabaseURL: url, Port: port}, nil
}
