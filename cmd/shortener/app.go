package main

import (
	"net/http"
)

type App struct {
	h http.Handler
}

func NewApp(handler http.Handler) *App {
	return &App{
		h: handler,
	}
}

func (app *App) Run(serverAddress string) {
	http.ListenAndServe(serverAddress, app.h)
}
