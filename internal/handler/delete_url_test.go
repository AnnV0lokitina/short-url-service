package handler

import (
	"context"
	"fmt"
	repoPkg "github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func ExampleHandler_DeleteBatch() {
	repo := repoPkg.NewMemoryRepo()
	service := service.NewService("http://localhost:8080", repo)

	ctx := context.TODO()

	go func() {
		service.CreateDeleteWorkerPull(ctx, 3)
	}()

	h := &Handler{
		Mux:     chi.NewMux(),
		service: service,
	}

	sendBody := strings.NewReader("fullURL")
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", sendBody)
	w := httptest.NewRecorder()
	h.SetURL().ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	cookies := res.Cookies()

	sendBody = strings.NewReader("fullURL1")
	req = httptest.NewRequest(http.MethodPost, "http://localhost:8080", sendBody)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w = httptest.NewRecorder()
	h.SetURL().ServeHTTP(w, req)

	// pause to make right order of requests
	time.Sleep(1 * time.Millisecond)

	sendBody = strings.NewReader("[\"27580c7e4c2c1de6435730c74bb9f8ca\",\"65f982c59e12c15c4d1633694a00f258\"]")
	req = httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/user/urls", sendBody)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w = httptest.NewRecorder()
	h.DeleteBatch().ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	fmt.Println(res.StatusCode)

	// pause to delete urls
	time.Sleep(500 * time.Millisecond)

	req = httptest.NewRequest(http.MethodGet, "http://localhost:8080/27580c7e4c2c1de6435730c74bb9f8ca", sendBody)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w = httptest.NewRecorder()
	h.GetURL().ServeHTTP(w, req)
	res = w.Result()

	defer res.Body.Close()
	fmt.Println(res.StatusCode)

	req = httptest.NewRequest(http.MethodGet, "http://localhost:8080/65f982c59e12c15c4d1633694a00f258", sendBody)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w = httptest.NewRecorder()
	h.GetURL().ServeHTTP(w, req)
	res = w.Result()

	defer res.Body.Close()
	fmt.Println(res.StatusCode)

	ctx.Done()

	time.Sleep(500 * time.Millisecond)

	// Output:
	// 202
	// 410
	// 410
}
