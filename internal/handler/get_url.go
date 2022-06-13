package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
)

// GetURL Follow by shorten url.
func (h *Handler) GetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		path := r.URL.Path
		checksum := path[1:]
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

// GetUserURLList Get list of urls, added by user.
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
