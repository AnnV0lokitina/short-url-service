package main

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"  envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	BataBaseDSN     string `env:"DATABASE_DSN" envDefault:""`
}

func initConfig() config {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
