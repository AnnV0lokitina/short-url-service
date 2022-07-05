package handler

import (
	"fmt"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/repo"
	"github.com/AnnV0lokitina/short-url-service/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleHandler_ShortenBatch() {
	repo := repoPkg.NewMemoryRepo()
	service := service.NewService("http://localhost:8080", repo)

	h := &Handler{
		Mux:     chi.NewMux(),
		service: service,
	}

	sendBody := strings.NewReader("[{\"correlation_id\":\"id\",\"original_url\":\"original url\"}," +
		"{\"correlation_id\":\"string id1\",\"original_url\":\"original url1\"}]")
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten/batch", sendBody)
	w := httptest.NewRecorder()
	h.ShortenBatch().ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	resBody, _ := io.ReadAll(res.Body)
	fmt.Println(string(resBody))

	// Output:
	// 201
	// [{"correlation_id":"id","short_url":"http://localhost:8080/d2a6fdf1db40a4efe500fa10cd71c939"},{"correlation_id":"string id1","short_url":"http://localhost:8080/0792df6a3cc8943351bbe3c338cae56a"}]
}
