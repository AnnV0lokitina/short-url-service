package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type App struct {
	h http.Handler
}

func NewApp(handler http.Handler) *App {
	return &App{
		h: handler,
	}
}

func (app *App) Run(ctx context.Context, serverAddress string) error {
	httpShutdownCh := make(chan struct{})
	server := &http.Server{Addr: serverAddress, Handler: app.h}

	go func() {
		<-ctx.Done()

		graceCtx, graceCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer graceCancel()

		if err := server.Shutdown(graceCtx); err != nil {
			log.Println(err)
		}
		httpShutdownCh <- struct{}{}
	}()

	err := server.ListenAndServe()
	<-httpShutdownCh
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
