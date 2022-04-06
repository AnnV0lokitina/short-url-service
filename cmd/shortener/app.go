package main

import (
	"context"
	"fmt"
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

		fmt.Println("graceCancel")

		if err := server.Shutdown(graceCtx); err != nil {
			fmt.Println("err")
			fmt.Println(err)
			log.Fatal(err)
		}
		fmt.Println(httpShutdownCh)
		httpShutdownCh <- struct{}{}
	}()

	err := server.ListenAndServe()
	<-httpShutdownCh
	fmt.Println("before ErrServerClosed")
	if err == http.ErrServerClosed {
		fmt.Println("ErrServerClosed")
		return nil
	}
	fmt.Println("NO ErrServerClosed")
	fmt.Println(err)
	return err
}
