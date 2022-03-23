package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

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

		// получение данных из запроса
		var itemInputList []JSONItemRequest
		if err := json.Unmarshal(request, &itemInputList); err != nil {
			http.Error(w, "Invalid request 7", http.StatusBadRequest)
			return
		}

		// конвертация в объекты приложения
		batchURLList := JSONListToURLList(itemInputList, h.BaseURL)
		// JSONItemRequest нельзя передавать в ядро бизнес логики
		// т.к. он лишь описывает формат полученных данных
		// в случае замены api (например на GRPC), JSONItemRequest будет бесполезен

		// запись данных приложения в репозиторий
		// функцию можно вынести в слой service, но она будет только проксировать вызов
		err = h.repo.AddBatch(ctx, userID, batchURLList)
		if err != nil {
			http.Error(w, "Error add batch", http.StatusBadRequest)
			return
		}

		// конвертация результата работы в JSON объекты для вывода пользователю
		outputList := URLListTOJSONList(batchURLList)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(&outputList); err != nil {
			http.Error(w, "Invalid request 9", http.StatusBadRequest)
			return
		}
	}
}
