package grpc

import (
	"context"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetUserURLList(ctx context.Context, in *pb.UserInfo) (*pb.AddrFullInfoListResponse, error) {
	userID, err := getUserID(in.GetUser())
	if err != nil {
		return &pb.AddrFullInfoListResponse{}, status.Error(codes.Internal, "Internal")
	}
	list, err := h.service.GetRepo().GetUserURLList(ctx, userID)
	if err != nil {
		return &pb.AddrFullInfoListResponse{}, status.Error(codes.NotFound, "No content")
	}
	resultList := make([]*pb.AddrFullInfoListItemResponse, 0, len(list))
	for _, item := range list {
		resultList = append(resultList, &pb.AddrFullInfoListItemResponse{
			ShortAddr:    item.Short,
			OriginalAddr: item.Original,
		})
	}
	return &pb.AddrFullInfoListResponse{
		List: resultList,
	}, nil
}
