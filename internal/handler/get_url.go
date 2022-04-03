package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		fmt.Println(shortURL)
		url, found, err := h.service.GetRepo().GetURL(ctx, shortURL)
		fmt.Println(found)
		if found {
			w.Header().Set("Location", url.Original)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		if err != nil {
			var labelErr *labelError.LabelError
			if errors.As(err, &labelErr) && labelErr.Label == "GONE" {
				fmt.Println("GONE")
				w.WriteHeader(http.StatusGone)
				return
			}
			http.Error(w, "Invalid request 4", http.StatusBadRequest)
			return
		}
		http.Error(w, "Not found 4", http.StatusBadRequest)
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
		list, found, err := h.service.GetRepo().GetUserURLList(ctx, userID)
		if err != nil || !found {
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
