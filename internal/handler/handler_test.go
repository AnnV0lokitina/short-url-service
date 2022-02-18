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

var tmpURL *entity.URL

type MockedRepo struct {
	mock.Mock
}

func (r *MockedRepo) SetURL(url *entity.URL) {
	tmpURL = url
}

func (r *MockedRepo) GetURL(checksum string) (*entity.URL, error) {
	if checksum == tmpURL.GetChecksum() {
		return tmpURL, nil
	}
	return nil, errors.New("no url saved")
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) *http.Response {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func TestHandler_ServeHTTP(t *testing.T) {
	repo := new(MockedRepo)
	handler := NewHandler(repo)
	url := entity.NewURLFromFullLink("fullURL")

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
				target: "/" + url.GetChecksum(),
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

	ts := httptest.NewServer(handler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := testRequest(t, ts, tt.request.method, tt.request.target, tt.request.body)
			assert.Equal(t, tt.result.code, resp.StatusCode)

			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			if tt.result.body != "" {
				assert.Equal(t, tt.result.body, string(respBody))
			}

			if tt.result.contentType != "" {
				assert.Equal(t, tt.result.contentType, resp.Header.Get("Content-Type"))
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	type args struct {
		repo Repo
	}

	repo := new(MockedRepo)

	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "create new handler",
			args: args{
				repo: repo,
			},
			want: &Handler{
				repo: repo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(tt.args.repo)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}
