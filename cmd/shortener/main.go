package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handlerPkg "github.com/AnnV0lokitina/short-url-service/internal/handler"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/repo"
	"github.com/AnnV0lokitina/short-url-service/internal/service"
	"github.com/AnnV0lokitina/short-url-service/internal/sqlrepo"
)

const nOfWorkers = 3

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	cfg := initConfig()
	initParams(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()
	repo, err := initRepo(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close(ctx)

	service := service.NewService(cfg.BaseURL, repo)
	handler := handlerPkg.NewHandler(service)
	application := NewApp(handler)

	go func() {
		service.CreateDeleteWorkerPull(ctx, nOfWorkers)
	}()

	err = application.Run(ctx, cfg.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}

func initRepo(ctx context.Context, cfg *config) (service.Repo, error) {
	if cfg.DataBaseDSN != "" {
		repository, err := sqlrepo.NewSQLRepo(ctx, cfg.DataBaseDSN)
		if err != nil {
			return nil, err
		}
		return repository, nil
	}
	if cfg.FileStoragePath != "" {
		repository, err := repoPkg.NewFileRepo(cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}
		return repository, nil
	}
	return repoPkg.NewMemoryRepo(), nil
}
