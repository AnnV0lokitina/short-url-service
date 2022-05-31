package handler

import (
	"fmt"
	repoPkg "github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleHandler_DeleteBatch() {
	repo := repoPkg.NewMemoryRepo()
	service := service.NewService("http://localhost:8080", repo)

	h := &Handler{
		Mux:     chi.NewMux(),
		service: service,
	}

	sendBody := strings.NewReader("[\"d2a6fdf1db40a4efe500fa10cd71c939\",\"0792df6a3cc8943351bbe3c338cae56a\"]")
	req := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/user/urls", sendBody)
	w := httptest.NewRecorder()
	h.DeleteBatch().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	fmt.Println(res.StatusCode)

	// Output:
	// 202
}
