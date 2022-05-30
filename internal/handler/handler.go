package handler

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"

	"github.com/AnnV0lokitina/short-url-service.git/internal/service"
)

const (
	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentType     = "Content-Type"
	encoding              = "gzip"
)

type Service interface {
	DeleteURLList(ctx context.Context, userID uint32, checksums []string) error
	GetRepo() service.Repo
	GetBaseURL() string
	SetBaseURL(baseURL string)
}

type Handler struct {
	*chi.Mux
	service Service
}

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

	h.Get("/debug/pprof/", pprof.Index)
	h.Get("/debug/pprof/cmdline", pprof.Cmdline)
	h.Get("/debug/pprof/profile", pprof.Profile)
	h.Get("/debug/pprof/symbol", pprof.Symbol)
	h.Get("/debug/pprof/trace", pprof.Trace)
	h.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	h.MethodNotAllowed(h.ExecIfNotAllowed())

	return h
}

func (h *Handler) ExecIfNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request 5", http.StatusBadRequest)
	}
}
