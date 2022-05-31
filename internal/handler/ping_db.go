package handler

import (
	"context"
	"net/http"
)

// PingDB Check if db is available.
func (h *Handler) PingDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		dbAvailable := h.service.GetRepo().PingBD(ctx)

		if dbAvailable {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}
}
