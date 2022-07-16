// Package shortener makes short url
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	handlerPkg "github.com/AnnV0lokitina/short-url-service/internal/handler"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/repo"
	servicePkg "github.com/AnnV0lokitina/short-url-service/internal/service"
	"github.com/AnnV0lokitina/short-url-service/internal/sqlrepo"
)

const nOfWorkers = 3

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	log.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	params := InitParams()
	cfg := InitConfig(params)
	log.Println(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		<-c
		cancel()
	}()
	repo, err := initRepo(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close(ctx)

	application, service := createApp(cfg, repo)
	go func() {
		service.CreateDeleteWorkerPull(ctx, nOfWorkers)
	}()

	err = application.Run(ctx, cfg.ServerAddress, cfg.EnableHTTPS)
	if err != nil {
		log.Fatal(err)
	}
}

func createApp(cfg *config, repo servicePkg.Repo) (*App, *servicePkg.Service) {
	service := servicePkg.NewService(cfg.BaseURL, repo, cfg.TrustedSubnet)
	handler := handlerPkg.NewHandler(service)
	return NewApp(handler), service
}

func initRepo(ctx context.Context, cfg *config) (servicePkg.Repo, error) {
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
