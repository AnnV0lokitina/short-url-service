package handler

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ServeHTTP(t *testing.T) {
	repository := repo.NewRepo()
	services := usecase.NewUsecase(repository)
	handler := NewHandler(services)

	fullURL := "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"

	url := entity.NewURL(fullURL, "")
	url.CreateUUID()

	type result struct {
		body           string
		headerLocation string
		code           int
	}
	tests := []struct {
		name    string
		request *http.Request
		result  result
	}{
		{
			name:    "test create incorrect method",
			request: httptest.NewRequest(http.MethodPut, "/", nil),
			result: result{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
			},
		},
		{
			name:    "test create url no body",
			request: httptest.NewRequest(http.MethodPost, "/", nil),
			result: result{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
			},
		},
		{
			name:    "test create url incorrect url",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("////%%%%%%")),
			result: result{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
			},
		},
		{
			name:    "test create url positive",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fullURL)),
			result: result{
				body:           url.GetShortURL(),
				headerLocation: "",
				code:           http.StatusCreated,
			},
		},
		{
			name:    "test read url positive",
			request: httptest.NewRequest(http.MethodGet, "/"+url.GetUUID(), nil),
			result: result{
				body:           "",
				headerLocation: url.GetFullURL(),
				code:           http.StatusTemporaryRedirect,
			},
		},
		{
			name:    "test read url no id",
			request: httptest.NewRequest(http.MethodGet, "/", nil),
			result: result{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, tt.request)
			res := w.Result()

			assert.Equal(t, res.StatusCode, tt.result.code)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if tt.result.body != "" && string(resBody) != tt.result.body {
				assert.Equal(t, string(resBody), tt.result.body)
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	type args struct {
		usecase *usecase.Usecase
	}

	repository := repo.NewRepo()
	services := usecase.NewUsecase(repository)

	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "create new handler",
			args: args{
				usecase: services,
			},
			want: &Handler{
				usecase: services,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(tt.args.usecase)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}
