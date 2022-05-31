package handler

import (
	"fmt"
	repoPkg "github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
)

func ExampleHandler_PingDB() {
	repo := repoPkg.NewMemoryRepo()
	service := service.NewService("http://localhost:8080", repo)

	h := &Handler{
		Mux:     chi.NewMux(),
		service: service,
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)
	w := httptest.NewRecorder()
	h.PingDB().ServeHTTP(w, req)

	fmt.Println(w.Result().StatusCode)
	w.Result().Body.Close()

	// Output:
	// 200
}
