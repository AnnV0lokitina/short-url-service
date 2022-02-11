package main

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"net/http"
)

type App struct {
	h *handler.Handler
}

func NewApp(handler *handler.Handler) *App {
	return &App{
		h: handler,
	}
}

func (app *App) Run() {
	http.Handle("/", app.h)

	http.ListenAndServe("localhost:8080", nil)
}
