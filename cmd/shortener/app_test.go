package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type AppInterface interface {
	Run(serverAddress string)
}

type MockedHandler struct {
	mock.Mock
}

func (h *MockedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("123"))
}

func TestApp_Run(t *testing.T) {
	type fields struct {
		h *MockedHandler
	}
	handler := new(MockedHandler)
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "run app",
			fields: fields{
				h: handler,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				h: tt.fields.h,
			}
			assert.Implements(t, (*AppInterface)(nil), app)
		})
	}
}

func TestNewApp(t *testing.T) {
	type args struct {
		handler *MockedHandler
	}

	handler := new(MockedHandler)
	tests := []struct {
		name string
		args args
		want *App
	}{
		{
			name: "test new app",
			args: args{
				handler: handler,
			},
			want: &App{
				h: handler,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewApp(tt.args.handler)
			assert.Equal(t, got, tt.want)
		})
	}
}
