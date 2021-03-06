package handler

import (
	"compress/gzip"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_repo"
	servicePkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_service"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, request testRequestStruct) *http.Response {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(request.method, request.target, request.body)
	if request.acceptEncoding != nil {
		req.Header.Set(headerAcceptEncoding, *request.acceptEncoding)
	}
	if request.cookie != nil {
		req.AddCookie(request.cookie)
	}

	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	if request.acceptEncoding == nil {
		resp.Header.Del(`Content-Encoding`)
	}
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
	service := new(servicePkg.MockedService)
	handler := NewHandler(service)
	ts := httptest.NewServer(handler)
	handler.service.SetBaseURL(ts.URL)

	tests := getTestsDataList(t, ts)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPkg.TmpURLError = tt.setURLError

			if tt.request.dbEnabled != nil {
				repoPkg.PingDB = *tt.request.dbEnabled
			}

			resp := testRequest(t, tt.request)
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
		service Service
		baseURL string
	}

	service := new(servicePkg.MockedService)
	service.SetBaseURL(cfg.BaseURL)

	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "create new handler",
			args: args{
				service: service,
				baseURL: cfg.BaseURL,
			},
			want: &Handler{
				service: service,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(tt.args.service)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}
