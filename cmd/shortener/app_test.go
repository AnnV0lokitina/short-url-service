package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedHandler struct {
	mock.Mock
}

func (h *MockedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("123"))
}

func TestNewApp(t *testing.T) {
	args := flag.Args()
	if len(args) == 0 || args[0] != "local" {
		t.Skip("Skipping testing in CI environment")
	}
	fContent, err := ioutil.ReadFile("defaults/defaults_run_test.json")
	assert.Nil(t, err)
	var cfg = config{}
	err = json.Unmarshal(fContent, &cfg)
	assert.Nil(t, err)

	handler := new(MockedHandler)
	app := NewApp(handler, nil)
	assert.IsType(t, &App{}, app)

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	shutdownCh := make(chan error)
	go func() {
		sendErr := app.Run(ctx, cfg.ServerAddress, false)
		shutdownCh <- sendErr
	}()
	gotErr := <-shutdownCh
	cancel()
	assert.Nil(t, gotErr)
}

func TestGetSelfSignedOrLetsEncryptCert(t *testing.T) {
	manager := &autocert.Manager{}
	f := getSelfSignedOrLetsEncryptCert(manager)
	assert.NotNil(t, f)
}

func TestCreateCertificate(t *testing.T) {
	hello := &tls.ClientHelloInfo{
		ServerName: "myhost.io",
	}
	_, err := createCertificate("certs", hello)
	assert.Nil(t, err)
}

func TestCreateCacheDir(t *testing.T) {
	manager := &autocert.Manager{}
	dir := createCacheDir(manager)
	var s string
	assert.IsType(t, s, dir)
}

func TestCreateServer(t *testing.T) {
	handler := new(MockedHandler)
	s := createServer(handler, "myhost.io", false)
	assert.Nil(t, s.TLSConfig)
	s1 := createServer(handler, "myhost.io", true)
	assert.NotNil(t, s1.TLSConfig)
}
