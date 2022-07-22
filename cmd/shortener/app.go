package main

import (
	"context"
	"crypto/tls"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

// App stores the handler.
type App struct {
	h http.Handler
	g *GRPCService
}

// NewApp create new App.
func NewApp(handler http.Handler, grpcService *GRPCService) *App {
	return &App{
		h: handler,
		g: grpcService,
	}
}

func createCertificate(dirCache string, hello *tls.ClientHelloInfo) (tls.Certificate, error) {
	keyFile := filepath.Join(dirCache, hello.ServerName+".key")
	crtFile := filepath.Join(dirCache, hello.ServerName+".crt")
	return tls.LoadX509KeyPair(crtFile, keyFile)
}

func createCacheDir(certManager *autocert.Manager) string {
	dirCache, ok := certManager.Cache.(autocert.DirCache)
	if ok {
		return string(dirCache)
	}
	return "certs"
}

func getSelfSignedOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache := createCacheDir(certManager)
		certificate, err := createCertificate(dirCache, hello)
		if err != nil {
			log.Printf("%s\nFalling back to Letsencrypt\n", err)
			return certManager.GetCertificate(hello)
		}
		log.Println("Loaded selfsigned certificate.")
		return &certificate, err
	}
}

func createServer(h http.Handler, serverAddress string, enableHTTPS bool) *http.Server {
	var server *http.Server
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
			Handler:   h,
			TLSConfig: tlsConfig,
		}
	} else {
		log.Println("http settings " + serverAddress)
		server = &http.Server{Addr: serverAddress, Handler: h}
	}
	return server
}

// Run Start the application.
func (app *App) Run(ctx context.Context, serverAddress string, enableHTTPS bool) error {
	var err error

	httpShutdownCh := make(chan struct{})
	server := createServer(app.h, serverAddress, enableHTTPS)

	if app.g != nil {
		go func() {
			listen, err := net.Listen("tcp", ":3200")
			if err != nil {
				log.Fatal(err)
			}

			pb.RegisterURLsServer(app.g.Server, app.g.Handler)
			log.Println("gRPC server starts")
			if err := app.g.Server.Serve(listen); err != nil {
				log.Fatal(err)
			}
		}()
	}

	go func() {
		<-ctx.Done()

		if app.g != nil {
			go func() {
				app.g.Server.GracefulStop()
				log.Println("stop grpc")
			}()
		}

		graceCtx, graceCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer graceCancel()

		if shutdownErr := server.Shutdown(graceCtx); shutdownErr != nil {
			log.Println(shutdownErr)
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
