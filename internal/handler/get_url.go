package handler

import (
	"context"
	"encoding/json"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) GetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		checksum := chi.URLParam(r, "id")
		shortURL := entity.CreateShortURL(checksum, h.BaseURL)
		url, found, err := h.repo.GetURL(ctx, shortURL)
		if !found || err != nil {
			http.Error(w, "Invalid request 4", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", url.Original)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *Handler) GetUserURLList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userID, err := authorization(w, r)
		if err != nil {
			http.Error(w, "Create user error", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		list, found, err := h.repo.GetUserURLList(ctx, userID)
		if !found || err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&list); err != nil {
			http.Error(w, "Error while json conversion", http.StatusBadRequest)
			return
		}
	}
}
