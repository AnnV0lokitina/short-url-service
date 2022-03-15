package handler

import (
	"encoding/json"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

type testRequestStruct struct {
	method         string
	target         string
	body           io.Reader
	acceptEncoding *string
}

type testResultStruct struct {
	body            string
	headerLocation  string
	code            int
	contentType     string
	contentEncoding *string
}

type testStruct struct {
	name        string
	setURLError bool
	request     testRequestStruct
	result      testResultStruct
}

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"  envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func newStringPtr(x string) *string {
	return &x
}

func createJSONEncodedResponse(t *testing.T, responseURL string) string {
	jsonResponse := JSONResponse{
		Result: responseURL,
	}

	jsonEncodedResponse, err := json.Marshal(jsonResponse)
	require.NoError(t, err)

	jsonEncodedResponse = append(jsonEncodedResponse, '\n')

	return string(jsonEncodedResponse)
}

func getTestsDataList(t *testing.T, cfg config) []testStruct {
	url := entity.NewURLFromFullLink("fullURL")
	return []testStruct{
		{
			name:        "test create incorrect method",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPut,
				target: "/",
				body:   nil,
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test create url no body",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/",
				body:   nil,
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test create url incorrect url",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/",
				body:   strings.NewReader("////%%%%%%"),
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test set url error",
			setURLError: true,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/",
				body:   strings.NewReader("fullURL"),
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test create url positive",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/",
				body:   strings.NewReader("fullURL"),
			},
			result: testResultStruct{
				body:           url.GetShortURL(cfg.BaseURL),
				headerLocation: "",
				code:           http.StatusCreated,
				contentType:    "text/plain; charset=UTF-8",
			},
		},
		{
			name:        "test create url with gzip positive",
			setURLError: false,
			request: testRequestStruct{
				method:         http.MethodPost,
				target:         "/",
				body:           strings.NewReader("fullURL"),
				acceptEncoding: newStringPtr(encoding),
			},
			result: testResultStruct{
				body:            url.GetShortURL(cfg.BaseURL),
				headerLocation:  "",
				code:            http.StatusCreated,
				contentType:     "text/plain; charset=UTF-8",
				contentEncoding: newStringPtr(encoding),
			},
		},
		{
			name:        "test json-api create incorrect method",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPut,
				target: "/api/shorten",
				body:   nil,
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test json-api create url no body",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/api/shorten",
				body:   nil,
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test json-api create url incorrect json",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/api/shorten",
				body:   strings.NewReader("{\"url:\"http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m\"}"),
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test json-api create url incorrect url",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/api/shorten",
				body:   strings.NewReader("{\"url\":\"////%%%%%%\"}"),
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test json-api set url error",
			setURLError: true,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/api/shorten",
				body:   strings.NewReader("{\"url\":\"fullURL\"}"),
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
		{
			name:        "test json-api create url positive",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPost,
				target: "/api/shorten",
				body:   strings.NewReader("{\"url\":\"fullURL\"}"),
			},
			result: testResultStruct{
				body:           createJSONEncodedResponse(t, url.GetShortURL(cfg.BaseURL)),
				headerLocation: "",
				code:           http.StatusCreated,
				contentType:    "application/json; charset=UTF-8",
			},
		},
		{
			name:        "test read url positive",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodGet,
				target: "/" + url.GetChecksum(),
				body:   nil,
			},
			result: testResultStruct{
				body:           "",
				headerLocation: url.GetFullURL(),
				code:           http.StatusTemporaryRedirect,
				contentType:    "",
			},
		},
		{
			name:        "test read url no id",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodGet,
				target: "/",
				body:   nil,
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusBadRequest,
				contentType:    "",
			},
		},
	}
}
