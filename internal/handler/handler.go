package handler

import (
	"context"
	"github.com/AnnV0lokitina/short-url-service.git/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentType     = "Content-Type"
	encoding              = "gzip"
)

// Service declare interface if services
type Service interface {
	DeleteURLList(ctx context.Context, userID uint32, checksums []string) error
	GetRepo() service.Repo
	GetBaseURL() string
	SetBaseURL(baseURL string)
}

// Handler keep information to handle requests.
type Handler struct {
	*chi.Mux
	service Service
}

// NewHandler Create new Handler.
func NewHandler(service Service) *Handler {
	h := &Handler{
		Mux:     chi.NewMux(),
		service: service,
	}

	h.Use(CompressMiddleware)

	h.Post("/", h.SetURL())
	h.Post("/api/shorten", h.SetURLFromJSON())
	h.Get("/{id}", h.GetURL())
	h.Get("/api/user/urls", h.GetUserURLList())
	h.Get("/ping", h.PingDB())
	h.Post("/api/shorten/batch", h.ShortenBatch())
	h.Delete("/api/user/urls", h.DeleteBatch())

	h.MethodNotAllowed(h.ExecIfNotAllowed())

	return h
}

// ExecIfNotAllowed Executed if url is bad.
func (h *Handler) ExecIfNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request 5", http.StatusBadRequest)
	}
}
