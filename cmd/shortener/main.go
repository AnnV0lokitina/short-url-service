package main

import (
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
)

func main() {
	repository := repo.NewRepo()
	handler := handlerPkg.NewHandler(repository)
	application := NewApp(handler)
	application.Run()
}
