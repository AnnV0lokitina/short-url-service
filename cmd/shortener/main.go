package main

import (
	"context"
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/sqlrepo"
	"log"
)

func main() {
	cfg := initConfig()
	initParams(cfg)

	ctx := context.Background()
	repo := initRepo(ctx, cfg)
	defer repo.Close(ctx)

	handler := handlerPkg.NewHandler(cfg.BaseURL, repo)
	application := NewApp(handler)
	application.Run(cfg.ServerAddress)
}

func initRepo(ctx context.Context, cfg config) handlerPkg.Repo {
	if cfg.BataBaseDSN != "" {
		repository, err := sqlrepo.NewSQLRepo(ctx, cfg.BataBaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		return repository
	}
	if cfg.FileStoragePath == "" {
		return repo.NewMemoryRepo()
	}
	repository, err := repo.NewFileRepo(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	return repository
}
