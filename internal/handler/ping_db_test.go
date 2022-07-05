package handler

import (
	"fmt"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/repo"
	"github.com/AnnV0lokitina/short-url-service/internal/service"
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
	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output:
	// 200
}
