package handler

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	netUrl "net/url"
)

type Usecase interface {
	SetURL(fullURL string) *entity.URL
	GetURL(uuid string) (*entity.URL, error)
}

type Handler struct {
	*chi.Mux
	usecase Usecase
}

func NewHandler(usecase Usecase) *Handler {
	h := &Handler{
		Mux:     chi.NewMux(),
		usecase: usecase,
	}

	h.Post("/", h.SetURL())
	h.Get("/{id}", h.GetURL())
	h.MethodNotAllowed(h.ExecIfNotAllowed())

	return h
}

func (h *Handler) ExecIfNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request 5", http.StatusBadRequest)
	}
}

func (h *Handler) SetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if err != nil || len(url) == 0 {
			http.Error(w, "Invalid request 1", http.StatusBadRequest)
			return
		}
		_, err = netUrl.Parse(string(url))
		if err != nil {
			http.Error(w, "Invalid request 2", http.StatusBadRequest)
			return
		}

		urlInfo := h.usecase.SetURL(string(url))

		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(urlInfo.GetShortURL()))
	}
}

func (h *Handler) GetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		urlInfo, err := h.usecase.GetURL(id)
		if err != nil {
			http.Error(w, "Invalid request 4", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", urlInfo.GetFullURL())
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
