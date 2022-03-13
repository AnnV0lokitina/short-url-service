package main

import (
	"flag"
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/caarlos0/env/v6"
	"log"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"  envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

func main() {
	var (
		cfg        config
		repository *repo.Repo
		err        error
	)

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File storage path")
	flag.Parse()

	if cfg.FileStoragePath == "" {
		repository = repo.NewMemoryRepo()
	} else {
		repository, err = repo.NewFileRepo(cfg.FileStoragePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer repository.Close()
	handler := handlerPkg.NewHandler(cfg.BaseURL, repository)
	application := NewApp(handler)
	application.Run(cfg.ServerAddress)
}
