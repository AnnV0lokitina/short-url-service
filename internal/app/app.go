package app

import (
	"net/http"
)

type App struct {
	h *Handler
}

func NewApp(handler *Handler) *App {
	return &App{
		h: handler,
	}
}

func (app *App) Run() {
	http.Handle("/", app.h)

	http.ListenAndServe("localhost:8080", nil)
}
