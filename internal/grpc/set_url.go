package grpc

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	netUrl "net/url"
)

// SetURL .
func (h *Handler) SetURL(ctx context.Context, in *pb.SetAddrRequest) (*pb.SetAddrResponse, error) {
	userID, err := getUserID(in.GetUser())
	if err != nil {
		return &pb.SetAddrResponse{}, status.Error(codes.Internal, "Internal")
	}
	_, err = netUrl.Parse(in.GetOriginalAddr())
	if err != nil {
		return &pb.SetAddrResponse{}, status.Error(codes.InvalidArgument, "Invalid request")
	}
	url := entity.NewURL(in.GetOriginalAddr(), h.service.GetBaseURL())
	err = h.service.GetRepo().SetURL(ctx, userID, url)
	if err != nil {
		var labelErr *labelError.LabelError
		if !errors.As(err, &labelErr) || labelErr.Label != labelError.TypeConflict {
			return &pb.SetAddrResponse{}, status.Error(codes.AlreadyExists, "Already exists")
		}
		return &pb.SetAddrResponse{}, status.Error(codes.InvalidArgument, "Invalid request")
	}

	return &pb.SetAddrResponse{
		User:   userID,
		Result: url.Short,
	}, nil
}
