package main

import (
	"context"
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/service"
	"github.com/AnnV0lokitina/short-url-service.git/internal/sqlrepo"
	"log"
)

func main() {
	cfg := initConfig()
	initParams(cfg)

	ctx := context.Background()
	repo := initRepo(ctx, cfg)
	defer repo.Close(ctx)

	service := service.NewService(cfg.BaseURL, repo)
	//service.ProcessDeleteRequests(ctx, runtime.NumCPU())
	handler := handlerPkg.NewHandler(service)

	application := NewApp(handler)
	application.Run(cfg.ServerAddress)
}

func initRepo(ctx context.Context, cfg *config) service.Repo {
	if cfg.DataBaseDSN != "" {
		repository, err := sqlrepo.NewSQLRepo(ctx, cfg.DataBaseDSN)
		if err != nil {
			log.Fatal(err)
		}
		return repository
	}
	if cfg.FileStoragePath != "" {
		repository, err := repo.NewFileRepo(cfg.FileStoragePath)
		if err != nil {
			log.Fatal(err)
		}
		return repository
	}
	return repo.NewMemoryRepo()
}
