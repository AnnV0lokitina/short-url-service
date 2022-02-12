package handler

import (
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var storedFullURL = ""

type MockedUsecase struct {
	mock.Mock
}

func (u *MockedUsecase) SetURL(fullURL string) *entity.URL {
	storedFullURL = fullURL
	return entity.NewURL(fullURL, "uuid")
}

func (u *MockedUsecase) GetURL(uuid string) (*entity.URL, error) {
	if uuid == "uuid" {
		url := entity.NewURL(storedFullURL, uuid)
		return url, nil
	}
	return nil, errors.New("no url saved")
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandler_ServeHTTP(t *testing.T) {
	service := new(MockedUsecase)
	handler := NewHandler(service)
	url := service.SetURL("fullURL")

	type request struct {
		method string
		target string
		body   io.Reader
	}
	type result struct {
		body           string
		headerLocation string
		code           int
		contentType    string
	}
	tests := []struct {
		name    string
		request request
		result  result
	}{
		{
			name: "test create incorrect method",
			request: request{
				method: http.MethodPut,
				target: "/",
				body:   nil,
			},
			result: result{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name: "test create url no body",
			request: request{
				method: http.MethodPost,
				target: "/",
				body:   nil,
			},
			result: result{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name: "test create url incorrect url",
			request: request{
				method: http.MethodPost,
				target: "/",
				body:   strings.NewReader("////%%%%%%"),
			},
			result: result{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name: "test create url positive",
			request: request{
				method: http.MethodPost,
				target: "/",
				body:   strings.NewReader("fullURL"),
			},
			result: result{
				body:           url.GetShortURL(),
				headerLocation: "",
				code:           http.StatusCreated,
				contentType:    "text/plain; charset=UTF-8",
			},
		},
		{
			name: "test read url positive",
			request: request{
				method: http.MethodGet,
				target: "/" + url.GetUUID(),
				body:   nil,
			},
			result: result{
				body:           "",
				headerLocation: url.GetFullURL(),
				code:           http.StatusTemporaryRedirect,
				contentType:    "",
			},
		},
		{
			name: "test read url no id",
			request: request{
				method: http.MethodGet,
				target: "/",
				body:   nil,
			},
			result: result{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(handler)
			defer ts.Close()

			resp, body := testRequest(t, ts, tt.request.method, tt.request.target, tt.request.body)
			assert.Equal(t, tt.result.code, resp.StatusCode)

			if tt.result.body != "" {
				assert.Equal(t, tt.result.body, body)
			}

			if tt.result.contentType != "" {
				assert.Equal(t, tt.result.contentType, resp.Header.Get("Content-Type"))
			}
		})
	}
}

//func TestHandler_ServeHTTP(t *testing.T) {
//	repository := repo.NewRepo()
//	services := usecase.NewUsecase(repository)
//	handler := NewHandler(services)
//
//	fullURL := "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
//
//	url := entity.NewURL(fullURL, "")
//	url.CreateUUID()
//
//	type result struct {
//		body           string
//		headerLocation string
//		code           int
//		contentType    string
//	}
//	tests := []struct {
//		name    string
//		request *http.Request
//		result  result
//	}{
//		{
//			name:    "test create incorrect method",
//			request: httptest.NewRequest(http.MethodPut, "/", nil),
//			result: result{
//				body:           "",
//				headerLocation: "",
//				code:           http.StatusBadRequest,
//				contentType:    "",
//			},
//		},
//		{
//			name:    "test create url no body",
//			request: httptest.NewRequest(http.MethodPost, "/", nil),
//			result: result{
//				body:           "",
//				headerLocation: "",
//				code:           http.StatusBadRequest,
//				contentType:    "",
//			},
//		},
//		{
//			name:    "test create url incorrect url",
//			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("////%%%%%%")),
//			result: result{
//				body:           "",
//				headerLocation: "",
//				code:           http.StatusBadRequest,
//				contentType:    "",
//			},
//		},
//		{
//			name:    "test create url positive",
//			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fullURL)),
//			result: result{
//				body:           url.GetShortURL(),
//				headerLocation: "",
//				code:           http.StatusCreated,
//				contentType:    "text/plain; charset=UTF-8",
//			},
//		},
//		{
//			name:    "test read url positive",
//			request: httptest.NewRequest(http.MethodGet, "/"+url.GetUUID(), nil),
//			result: result{
//				body:           "",
//				headerLocation: url.GetFullURL(),
//				code:           http.StatusTemporaryRedirect,
//				contentType:    "",
//			},
//		},
//		{
//			name:    "test read url no id",
//			request: httptest.NewRequest(http.MethodGet, "/", nil),
//			result: result{
//				body:           "",
//				headerLocation: "",
//				code:           http.StatusBadRequest,
//				contentType:    "",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			w := httptest.NewRecorder()
//			handler.ServeHTTP(w, tt.request)
//			res := w.Result()
//
//			assert.Equal(t, tt.result.code, res.StatusCode)
//
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//			require.NoError(t, err)
//
//			if tt.result.body != "" {
//				assert.Equal(t, tt.result.body, string(resBody))
//			}
//
//			if tt.result.contentType != "" {
//				assert.Equal(t, tt.result.contentType, res.Header.Get("Content-Type"))
//			}
//		})
//	}
//}

func TestNewHandler(t *testing.T) {
	type args struct {
		usecase Usecase
	}

	services := new(MockedUsecase)

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
