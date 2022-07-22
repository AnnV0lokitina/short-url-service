package grpc

import (
	"context"
	"errors"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	headerIP = "X-Real-IP"
)

func (h *Handler) GetStats(ctx context.Context, _ *pb.StatRequest) (*pb.StatsResponse, error) {
	var ipStr string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get(headerIP)
		if len(values) > 0 {
			ipStr = values[0]
		}
	}
	stats, err := h.service.GetStats(ctx, ipStr)
	if err != nil {
		var labelErr *labelError.LabelError
		if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeForbidden {
			return &pb.StatsResponse{}, status.Error(codes.PermissionDenied, "Forbidden")
		}
		return &pb.StatsResponse{}, status.Error(codes.InvalidArgument, "Invalid request")
	}
	return &pb.StatsResponse{
		NAddr:  uint32(stats.URLs),
		NUsers: uint32(stats.Users),
	}, nil
}
