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

func (app *App) Run() {
	http.ListenAndServe("localhost:8080", app.h)
}
