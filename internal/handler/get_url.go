package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service.git/pkg/error"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) GetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		checksum := chi.URLParam(r, "id")
		shortURL := entity.CreateShortURL(checksum, h.service.GetBaseURL())
		url, err := h.service.GetRepo().GetURL(ctx, shortURL)
		if err == nil {
			w.Header().Set("Location", url.Original)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		var labelErr *labelError.LabelError
		if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeGone {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.Error(w, "Invalid request 4", http.StatusBadRequest)
	}
}

func (h *Handler) GetUserURLList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userID, err := authorizeUserAndSetCookie(w, r)
		if err != nil {
			http.Error(w, "Create user error", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		list, err := h.service.GetRepo().GetUserURLList(ctx, userID)
		if err != nil {
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
