package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service.git/pkg/error"
	"io"
	"net/http"
	netUrl "net/url"
)

type JSONRequest struct {
	URL string `json:"url"`
}

type JSONResponse struct {
	Result string `json:"result"`
}

func (h *Handler) SetURLFromJSON() http.HandlerFunc {
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

		var parsedRequest JSONRequest
		if err := json.Unmarshal(request, &parsedRequest); err != nil {
			http.Error(w, "Invalid request 7", http.StatusBadRequest)
			return
		}

		_, err = netUrl.Parse(parsedRequest.URL)
		if err != nil {
			http.Error(w, "Invalid request 8", http.StatusBadRequest)
			return
		}

		url := entity.NewURL(parsedRequest.URL, h.BaseURL)
		err = h.repo.SetURL(ctx, userID, url)
		if err != nil {
			var labelErr *labelError.LabelError
			if !errors.As(err, &labelErr) || labelErr.Label != "CONFLICT" {
				http.Error(w, "Invalid request 10", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusConflict)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
		}

		jsonResponse := JSONResponse{
			Result: url.Short,
		}
		if err := json.NewEncoder(w).Encode(&jsonResponse); err != nil {
			http.Error(w, "Invalid request 9", http.StatusBadRequest)
			return
		}
	}
}

func (h *Handler) SetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userID, err := authorization(w, r)
		if err != nil {
			http.Error(w, "Create user error", http.StatusBadRequest)
			return
		}
		url, err := io.ReadAll(r.Body)
		if err != nil || len(url) == 0 {
			http.Error(w, "Invalid request 1", http.StatusBadRequest)
			return
		}
		_, err = netUrl.Parse(string(url))
		if err != nil {
			http.Error(w, "Invalid request 2", http.StatusBadRequest)
			return
		}

		urlPair := entity.NewURL(string(url), h.BaseURL)
		err = h.repo.SetURL(ctx, userID, urlPair)
		if err != nil {
			var labelErr *labelError.LabelError
			if !errors.As(err, &labelErr) || labelErr.Label != "CONFLICT" {
				http.Error(w, "Invalid request 10", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			w.WriteHeader(http.StatusConflict)
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
		}
		w.Write([]byte(urlPair.Short))
	}
}
