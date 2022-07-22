package grpc

import (
	"context"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) SetShortenBatch(ctx context.Context, in *pb.AddrListRequest) (*pb.AddrListResponse, error) {
	userID, err := getUserID(in.GetUser())
	if err != nil {
		return &pb.AddrListResponse{}, status.Error(codes.Internal, "Internal")
	}
	list := make([]*entity.BatchURLItem, 0, len(in.GetList()))
	for _, item := range in.GetList() {
		urlItem := entity.NewBatchURLItem(
			item.CorrelationId,
			item.OriginalAddr,
			h.service.GetBaseURL(),
		)
		list = append(list, urlItem)
	}
	err = h.service.GetRepo().AddBatch(ctx, userID, list)
	if err != nil {
		return &pb.AddrListResponse{}, status.Error(codes.InvalidArgument, "Invalid request")
	}
	outputList := make([]*pb.AddrListItemResponse, 0, len(list))
	for _, item := range list {
		i := &pb.AddrListItemResponse{
			CorrelationId: item.CorrelationID,
			ShortAddr:     item.URL.Short,
		}
		outputList = append(outputList, i)
	}
	return &pb.AddrListResponse{
		User: userID,
		List: outputList,
	}, nil
}
