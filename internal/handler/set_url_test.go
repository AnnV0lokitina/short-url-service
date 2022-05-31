package handler

import (
	"fmt"
	repoPkg "github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleHandler_SetURL() {
	repo := repoPkg.NewMemoryRepo()
	service := service.NewService("http://localhost:8080", repo)

	h := &Handler{
		Mux:     chi.NewMux(),
		service: service,
	}

	sendBody := strings.NewReader("fullURL")
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", sendBody)
	w := httptest.NewRecorder()
	h.SetURL().ServeHTTP(w, req)

	fmt.Println(w.Result().StatusCode)
	resBody, _ := io.ReadAll(w.Result().Body)
	fmt.Println(string(resBody))
	w.Result().Body.Close()

	// Output:
	// 201
	// http://localhost:8080/27580c7e4c2c1de6435730c74bb9f8ca
}

func ExampleHandler_SetURLFromJSON() {
	repo := repoPkg.NewMemoryRepo()
	service := service.NewService("http://localhost:8080", repo)

	h := &Handler{
		Mux:     chi.NewMux(),
		service: service,
	}

	sendBody := strings.NewReader("{\"url\": \"fullURL\"}")
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", sendBody)
	w := httptest.NewRecorder()
	h.SetURLFromJSON().ServeHTTP(w, req)

	fmt.Println(w.Result().StatusCode)
	resBody, _ := io.ReadAll(w.Result().Body)
	fmt.Println(string(resBody))
	w.Result().Body.Close()

	// Output:
	// 201
	// {"result":"http://localhost:8080/27580c7e4c2c1de6435730c74bb9f8ca"}
}
