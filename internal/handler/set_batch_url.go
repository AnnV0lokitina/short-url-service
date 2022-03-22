package handler

import (
	"context"
	"encoding/json"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"io"
	"net/http"
)

type itemInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type itemOutput struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (ii itemInput) toBatchURLItem(serverAddress string) *entity.BatchURLItem {
	return entity.NewBatchURLItem(
		ii.CorrelationID,
		ii.OriginalURL,
		serverAddress,
	)
}

func (io *itemOutput) fromBatchURLItem(item *entity.BatchURLItem) {
	io.CorrelationID = item.CorrelationID
	io.ShortURL = item.URL.Short
}

func (h *Handler) ShortenBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		userID, err := authorization(w, r)
		if err != nil {
			http.Error(w, "Create user error", http.StatusBadRequest)
			return
		}

		request, err := io.ReadAll(r.Body)
		if err != nil || len(request) == 0 {
			http.Error(w, "Invalid request 6", http.StatusBadRequest)
			return
		}

		var itemInputList []itemInput
		if err := json.Unmarshal(request, &itemInputList); err != nil {
			http.Error(w, "Invalid request 7", http.StatusBadRequest)
			return
		}

		list := make([]*entity.BatchURLItem, 0)
		for _, item := range itemInputList {
			urlItem := item.toBatchURLItem(h.BaseURL)
			list = append(list, urlItem)
		}

		err = h.repo.AddBatch(ctx, userID, list)
		if err != nil {
			http.Error(w, "Error add batch", http.StatusBadRequest)
			return
		}

		outputList := make([]itemOutput, 0)
		for _, item := range list {
			i := &itemOutput{}
			i.fromBatchURLItem(item)
			outputList = append(outputList, *i)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(&outputList); err != nil {
			http.Error(w, "Invalid request 9", http.StatusBadRequest)
			return
		}
	}
}
