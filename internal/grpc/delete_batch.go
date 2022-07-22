package grpc

import (
	"context"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) DeleteBatch(ctx context.Context, in *pb.DeleteBatchRequest) (*pb.TextResponse, error) {
	userID, err := getUserID(in.GetUser())
	if err != nil {
		return &pb.TextResponse{}, status.Error(codes.Internal, "Internal")
	}
	err = h.service.DeleteURLList(ctx, userID, in.GetChecksumList())
	if err != nil {
		return &pb.TextResponse{}, status.Error(codes.InvalidArgument, "Invalid request")
	}

	return &pb.TextResponse{
		Result: "OK",
	}, nil
}
