package grpc

import (
	"context"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	"github.com/AnnV0lokitina/short-url-service/internal/service"
	"github.com/AnnV0lokitina/short-url-service/pkg/userid"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
)

// Service interface describes business-logic layer
type Service interface {
	DeleteURLList(ctx context.Context, userID uint32, checksums []string) error
	GetRepo() service.Repo
	GetBaseURL() string
	SetBaseURL(baseURL string)
	GetStats(ctx context.Context, ipStr string) (entity.Stats, error)
}

// Handler structure holds dependencies for server handlers.
type Handler struct {
	pb.UnimplementedURLsServer

	service Service
}

func NewHandler(service Service) *Handler {
	h := &Handler{}
	h.service = service
	return h
}

func getUserID(id uint32) (uint32, error) {
	if id > 0 {
		return id, nil
	}
	return userid.GenerateUserID()
}
