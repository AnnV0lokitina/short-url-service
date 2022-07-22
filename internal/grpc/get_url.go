package grpc

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetURL(ctx context.Context, in *pb.GetAddrRequest) (*pb.GetAddrResponse, error) {
	shortURL := entity.CreateShortURL(in.GetChecksum(), h.service.GetBaseURL())
	url, err := h.service.GetRepo().GetURL(ctx, shortURL)
	if err != nil {
		var labelErr *labelError.LabelError
		if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeGone {
			return &pb.GetAddrResponse{}, status.Error(codes.FailedPrecondition, "URL deleted")
		}
		return &pb.GetAddrResponse{}, status.Error(codes.InvalidArgument, "Invalid request")
	}
	return &pb.GetAddrResponse{
		OriginalAddr: url.Original,
	}, nil
}
