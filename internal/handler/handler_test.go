package handler

import (
	"compress/gzip"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	tmpURL      *entity.URL
	tmpURLError = false
)

type MockedRepo struct {
	mock.Mock
}

func (r *MockedRepo) SetURL(url *entity.URL) error {
	if tmpURLError == true {
		return errors.New("error")
	}
	tmpURL = url
	return nil
}

func (r *MockedRepo) GetURL(checksum string) (*entity.URL, error) {
	if checksum == tmpURL.GetChecksum() {
		return tmpURL, nil
	}
	return nil, errors.New("no url saved")
}

func testRequest(t *testing.T, ts *httptest.Server, request testRequestStruct) *http.Response {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(request.method, ts.URL+request.target, request.body)
	if request.acceptEncoding != nil {
		req.Header.Set(headerAcceptEncoding, *request.acceptEncoding)
	}

	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func getResponseReader(t *testing.T, resp *http.Response) io.Reader {
	var reader io.Reader
	if resp.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		reader = gz
		defer gz.Close()
	} else {
		reader = resp.Body
	}

	return reader
}

func TestHandler_ServeHTTP(t *testing.T) {
	cfg := config{}
	err := env.Parse(&cfg)
	require.NoError(t, err)

	repo := new(MockedRepo)
	handler := NewHandler(cfg.BaseURL, repo)
	tests := getTestsDataList(t, cfg)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpURLError = tt.setURLError
			resp := testRequest(t, ts, tt.request)
			assert.Equal(t, tt.result.code, resp.StatusCode)

			if tt.result.contentEncoding != nil {
				assert.Equal(t, *tt.result.contentEncoding, resp.Header.Get(headerContentEncoding))
			}

			defer resp.Body.Close()

			reader := getResponseReader(t, resp)

			respBody, err := ioutil.ReadAll(reader)
			require.NoError(t, err)

			if tt.result.body != "" {
				assert.Equal(t, tt.result.body, string(respBody))
			}

			if tt.result.contentType != "" {
				assert.Equal(t, tt.result.contentType, resp.Header.Get(headerContentType))
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	cfg := config{}
	err := env.Parse(&cfg)
	require.NoError(t, err)

	type args struct {
		repo    Repo
		baseURL string
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
				repo:    repo,
				baseURL: cfg.BaseURL,
			},
			want: &Handler{
				repo:    repo,
				baseURL: cfg.BaseURL,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(tt.args.baseURL, tt.args.repo)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}
