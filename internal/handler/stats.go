package handler

import (
	"context"
	"encoding/json"
	"errors"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
	"net/http"
)

// GetStats Gets statistic of urls and users.
func (h *Handler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ipStr := r.Header.Get("X-Real-IP")
		stats, err := h.service.GetStats(ctx, ipStr)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			var labelErr *labelError.LabelError
			if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeForbidden {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&stats); err != nil {
			http.Error(w, "Error while json conversion", http.StatusBadRequest)
			return
		}
	}
}
