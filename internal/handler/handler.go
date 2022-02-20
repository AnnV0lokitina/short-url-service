package handler

import (
	"encoding/json"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	netUrl "net/url"
)

type Repo interface {
	SetURL(url *entity.URL)
	GetURL(checksum string) (*entity.URL, error)
}

type Handler struct {
	*chi.Mux
	repo    Repo
	baseURL string
}

type JSONRequest struct {
	URL string `json:"url"`
}

type JSONResponse struct {
	Result string `json:"result"`
}

func NewHandler(baseURL string, repo Repo) *Handler {
	h := &Handler{
		Mux:     chi.NewMux(),
		repo:    repo,
		baseURL: baseURL,
	}

	h.Post("/", h.SetURL())
	h.Post("/api/shorten", h.SetURLFromJSON())
	h.Get("/{id}", h.GetURL())
	h.MethodNotAllowed(h.ExecIfNotAllowed())

	return h
}

func (h *Handler) ExecIfNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request 5", http.StatusBadRequest)
	}
}

func (h *Handler) SetURLFromJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := io.ReadAll(r.Body)
		if err != nil || len(request) == 0 {
			http.Error(w, "Invalid request 6", http.StatusBadRequest)
			return
		}

		var parsedRequest JSONRequest
		if err := json.Unmarshal(request, &parsedRequest); err != nil {
			http.Error(w, "Invalid request 7", http.StatusBadRequest)
			return
		}

		_, err = netUrl.Parse(parsedRequest.URL)
		if err != nil {
			http.Error(w, "Invalid request 8", http.StatusBadRequest)
			return
		}

		urlInfo := entity.NewURLFromFullLink(parsedRequest.URL)
		h.repo.SetURL(urlInfo)

		jsonResponse := JSONResponse{
			Result: urlInfo.GetShortURL(h.baseURL),
		}

		response, err := json.Marshal(jsonResponse)
		if err != nil {
			http.Error(w, "Invalid request 9", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
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

		urlInfo := entity.NewURLFromFullLink(string(url))
		h.repo.SetURL(urlInfo)

		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(urlInfo.GetShortURL(h.baseURL)))
	}
}

func (h *Handler) GetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		urlInfo, err := h.repo.GetURL(id)
		if err != nil {
			http.Error(w, "Invalid request 4", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", urlInfo.GetFullURL())
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
