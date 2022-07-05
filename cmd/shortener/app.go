package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

// App stores the handler.
type App struct {
	h http.Handler
}

// NewApp create new App.
func NewApp(handler http.Handler) *App {
	return &App{
		h: handler,
	}
}

func getSelfSignedOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certManager.Cache.(autocert.DirCache)
		if !ok {
			dirCache = "certs"
		}

		keyFile := filepath.Join(string(dirCache), hello.ServerName+".key")
		crtFile := filepath.Join(string(dirCache), hello.ServerName+".crt")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			log.Printf("%s\nFalling back to Letsencrypt\n", err)
			return certManager.GetCertificate(hello)
		}
		log.Println("Loaded selfsigned certificate.")
		return &certificate, err
	}
}

// Run Start the application.
func (app *App) Run(ctx context.Context, serverAddress string, enableHTTPS bool) error {
	var err error
	var server *http.Server

	httpShutdownCh := make(chan struct{})

	if enableHTTPS {
		log.Println("https settings")
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("certs"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("myhost.io"),
		}
		tlsConfig := manager.TLSConfig()
		tlsConfig.MaxVersion = tls.VersionTLS10
		tlsConfig.GetCertificate = getSelfSignedOrLetsEncryptCert(manager)
		server = &http.Server{
			Addr:      ":8081",
			Handler:   app.h,
			TLSConfig: tlsConfig,
		}
	} else {
		log.Println("http settings " + serverAddress)
		server = &http.Server{Addr: serverAddress, Handler: app.h}
	}

	go func() {
		<-ctx.Done()

		graceCtx, graceCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer graceCancel()

		if err = server.Shutdown(graceCtx); err != nil {
			log.Println(err)
		}
		httpShutdownCh <- struct{}{}
	}()

	if enableHTTPS {
		log.Println("start https")
		err = server.ListenAndServeTLS("", "")
	} else {
		log.Println("start http")
		err = server.ListenAndServe()
	}

	<-httpShutdownCh
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
