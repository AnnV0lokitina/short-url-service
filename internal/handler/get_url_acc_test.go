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
	"time"
)

func ExampleHandler_GetURL() {
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

	req = httptest.NewRequest(http.MethodGet, "http://localhost:8080/27580c7e4c2c1de6435730c74bb9f8ca", sendBody)
	w = httptest.NewRecorder()
	h.GetURL().ServeHTTP(w, req)
	res := w.Result()

	defer res.Body.Close()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header.Get("Location"))

	// Output:
	// 307
	// fullURL
}

func ExampleHandler_GetUserURLList() {
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
	res := w.Result()
	defer res.Body.Close()
	cookies := res.Cookies()

	time.Sleep(10 * time.Millisecond)

	sendBody = strings.NewReader("fullURL1")
	req = httptest.NewRequest(http.MethodPost, "http://localhost:8080", sendBody)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w = httptest.NewRecorder()
	h.SetURL().ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/urls", sendBody)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w = httptest.NewRecorder()
	h.GetUserURLList().ServeHTTP(w, req)
	resCheck := w.Result()
	defer resCheck.Body.Close()

	fmt.Println(resCheck.StatusCode)
	resBody, _ := io.ReadAll(resCheck.Body)
	fmt.Println(string(resBody))

	// Output:
	// 200
	// [{"short_url":"http://localhost:8080/27580c7e4c2c1de6435730c74bb9f8ca","original_url":"fullURL"},{"short_url":"http://localhost:8080/65f982c59e12c15c4d1633694a00f258","original_url":"fullURL1"}]
}
