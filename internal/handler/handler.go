package handler

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentType     = "Content-Type"
	encoding              = "gzip"
)

type Repo interface {
	SetURL(userID uint32, url *entity.URL) error
	GetURL(checksum string) (*entity.URL, error)
	GetUserURLList(id uint32) ([]*entity.URL, bool)
}

type Handler struct {
	*chi.Mux
	repo    Repo
	BaseURL string
}

func NewHandler(baseURL string, repo Repo) *Handler {
	h := &Handler{
		Mux:     chi.NewMux(),
		repo:    repo,
		BaseURL: baseURL,
	}

	h.Use(CompressMiddleware)

	h.Post("/", h.SetURL())
	h.Post("/api/shorten", h.SetURLFromJSON())
	h.Get("/{id}", h.GetURL())
	h.Get("/api/user/urls", h.GetUserURLList())
	h.MethodNotAllowed(h.ExecIfNotAllowed())

	return h
}

func (h *Handler) ExecIfNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request 5", http.StatusBadRequest)
	}
}
