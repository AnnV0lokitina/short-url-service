package main

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/app"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/usecase"
	//"github.com/AnnV0lokitina/short-url-service/internal/app"
)

func main() {
	repository := repo.NewRepo()
	usecase := usecase.NewUsecase(repository)
	handler := app.NewHandler(usecase)
	application := app.NewApp(handler)
	application.Run()
}
