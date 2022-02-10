package app

import (
	handlerPkg "github.com/AnnV0lokitina/short-url-service.git/internal/handler"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/AnnV0lokitina/short-url-service.git/internal/usecase"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp_Run(t *testing.T) {
	type fields struct {
		h *handlerPkg.Handler
	}
	repository := repo.NewRepo()
	services := usecase.NewUsecase(repository)
	handler := handlerPkg.NewHandler(services)
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
			app.Run()
		})
	}
}

func TestNewApp(t *testing.T) {
	type args struct {
		handler *handlerPkg.Handler
	}
	repository := repo.NewRepo()
	services := usecase.NewUsecase(repository)
	handler := handlerPkg.NewHandler(services)
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
