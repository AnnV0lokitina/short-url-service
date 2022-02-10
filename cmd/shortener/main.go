package main

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/app"
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/usecase"
)

func main() {
	repository := repo.NewRepo()
	services := usecase.NewUsecase(repository)
	handler := handlerPkg.NewHandler(services)
	application := app.NewApp(handler)
	application.Run()
}
