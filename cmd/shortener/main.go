package main

import (
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/usecase"
)

func main() {
	repository := repo.NewRepo()
	services := usecase.NewUsecase(repository)
	handler := handlerPkg.NewHandler(services)
	application := NewApp(handler)
	application.Run()
}
