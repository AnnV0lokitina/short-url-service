package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) DeleteBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userID, err := authorizeUserAndSetCookie(w, r)
		if err != nil {
			http.Error(w, "Create user error", http.StatusBadRequest)
			return
		}

		request, err := io.ReadAll(r.Body)
		if err != nil || len(request) == 0 {
			http.Error(w, "Invalid request 6", http.StatusBadRequest)
			return
		}

		// получение данных из запроса
		var itemChecksumList []string
		if err := json.Unmarshal(request, &itemChecksumList); err != nil {
			http.Error(w, "Invalid request 7", http.StatusBadRequest)
			return
		}

		err = h.service.DeleteURLList(ctx, userID, itemChecksumList)
		if err != nil {
			http.Error(w, "Create user error 8", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
