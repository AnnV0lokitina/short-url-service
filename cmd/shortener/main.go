package main

import (
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress   string  `env:"SERVER_ADDRESS"  envDefault:"localhost:8080"`
	BaseURL         string  `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	repository, err := repo.NewRepo(cfg.FileStoragePath)
	if err != nil {
		panic(err)
	}
	defer repository.Close()
	handler := handlerPkg.NewHandler(cfg.BaseURL, repository)
	application := NewApp(handler)
	application.Run(cfg.ServerAddress)
}
