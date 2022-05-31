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

	defer w.Result().Body.Close()
	fmt.Println(w.Result().StatusCode)
	fmt.Println(w.Result().Header.Get("Location"))

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
	defer w.Result().Body.Close()
	cookies := w.Result().Cookies()

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
	w.Result().Body.Close()

	fmt.Println(w.Result().StatusCode)
	resBody, _ := io.ReadAll(w.Result().Body)
	fmt.Println(string(resBody))
	w.Result().Body.Close()

	// Output:
	// 200
	// [{"short_url":"http://localhost:8080/27580c7e4c2c1de6435730c74bb9f8ca","original_url":"fullURL"},{"short_url":"http://localhost:8080/65f982c59e12c15c4d1633694a00f258","original_url":"fullURL1"}]
}
