package handler

import (
	"encoding/json"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testRequestStruct struct {
	method         string
	target         string
	body           io.Reader
	acceptEncoding *string
	cookie         *http.Cookie
	dbEnabled      *bool
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

func newBoolPtr(x bool) *bool {
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

func getTestsDataList(t *testing.T, ts *httptest.Server) []testStruct {
	url := entity.NewURL("fullURL", ts.URL)
	return []testStruct{
		{
			name:        "test create incorrect method",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodPut,
				target: ts.URL + "/",
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
				target: ts.URL + "/",
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
				target: ts.URL + "/",
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
				target: ts.URL + "/",
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
				target: ts.URL + "/",
				body:   strings.NewReader("fullURL"),
			},
			result: testResultStruct{
				body:           url.Short,
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
				target:         ts.URL + "/",
				body:           strings.NewReader("fullURL"),
				acceptEncoding: newStringPtr(encoding),
			},
			result: testResultStruct{
				body:            url.Short,
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
				target: ts.URL + "/api/shorten",
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
				target: ts.URL + "/api/shorten",
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
				target: ts.URL + "/api/shorten",
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
				target: ts.URL + "/api/shorten",
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
				target: ts.URL + "/api/shorten",
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
				target: ts.URL + "/api/shorten",
				body:   strings.NewReader("{\"url\":\"fullURL\"}"),
			},
			result: testResultStruct{
				body:           createJSONEncodedResponse(t, url.Short),
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
				target: url.Short,
				body:   nil,
			},
			result: testResultStruct{
				body:           "",
				headerLocation: url.Original,
				code:           http.StatusTemporaryRedirect,
				contentType:    "",
			},
		},
		{
			name:        "test read url no id",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/",
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
			name:        "test get user url nagative",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/urls",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "login",
					Value:    generateLogin(4444),
					HttpOnly: false,
				},
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusNoContent,
				contentType:    "application/json; charset=UTF-8",
			},
		},
		{
			name:        "test get user url positive",
			setURLError: false,
			request: testRequestStruct{
				method: http.MethodGet,
				target: ts.URL + "/api/user/urls",
				body:   nil,
				cookie: &http.Cookie{
					Name:     "login",
					Value:    generateLogin(1234),
					HttpOnly: false,
				},
				acceptEncoding: nil,
			},
			result: testResultStruct{
				body:           "[{\"short_url\":\"short\",\"original_url\":\"original\"}]\n",
				headerLocation: "",
				code:           http.StatusOK,
				contentType:    "application/json; charset=UTF-8",
			},
		},
		{
			name:        "test ping nagative",
			setURLError: false,
			request: testRequestStruct{
				method:    http.MethodGet,
				target:    ts.URL + "/ping",
				body:      nil,
				dbEnabled: newBoolPtr(false),
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusInternalServerError,
				contentType:    "",
			},
		},
		{
			name:        "test ping positive",
			setURLError: false,
			request: testRequestStruct{
				method:    http.MethodGet,
				target:    ts.URL + "/ping",
				body:      nil,
				dbEnabled: newBoolPtr(true),
			},
			result: testResultStruct{
				body:           "",
				headerLocation: "",
				code:           http.StatusOK,
				contentType:    "",
			},
		},
	}
}
