package main

import (
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"  envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	repository := repo.NewRepo()
	handler := handlerPkg.NewHandler(cfg.BaseURL, repository)
	application := NewApp(handler)
	application.Run(cfg.ServerAddress)
}
