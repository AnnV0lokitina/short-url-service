package handler

import (
	"io"
	"net/http"
	netUrl "net/url"

	"github.com/AnnV0lokitina/short-url-service.git/internal/usecase"
)

type Handler struct {
	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
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

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(urlInfo.GetShortURL()))
	case http.MethodGet:
		path := r.URL.Path
		id := path[1:]
		urlInfo, err := h.usecase.GetURL(id)
		if err != nil {
			http.Error(w, "Invalid request 4", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", urlInfo.GetFullURL())
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		http.Error(w, "Invalid request 5", http.StatusBadRequest)
	}
}
