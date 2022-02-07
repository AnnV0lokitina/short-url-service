package app

import (
	"fmt"
	"net/http"

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
		if err := r.ParseForm(); err != nil {
			http.Error(w, "No params", 400)
			return
		}
		url := r.FormValue("url")
		if url == "" {
			http.Error(w, "No url", 400)
			return
		}

		uuid, shortURL := h.usecase.SetURL(url)

		showInfo := fmt.Sprintf("%s: %s", uuid, shortURL)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(showInfo))
	case http.MethodGet:
		path := r.URL.Path
		id := path[1:]
		full, short, err := h.usecase.GetURL(id)
		if err != nil {
			http.Error(w, "Invalid request", 400)
		}
		w.Header().Set("Location", full)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(short))
	default:
		http.Error(w, "Invalid request", 400)
	}
}
