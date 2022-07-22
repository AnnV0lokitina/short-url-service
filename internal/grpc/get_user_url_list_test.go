package grpc

import (
	"context"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_repo"
	servicePkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_service"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserURLList(t *testing.T) {
	service := servicePkg.NewMockedService()
	service.SetBaseURL("baseURL")
	h := NewHandler(service)
	ctx := context.TODO()

	in := &pb.UserInfo{
		User: repoPkg.RightUser,
	}
	urlList := []*entity.URL{
		entity.NewURL("original1", "baseURL"),
		entity.NewURL("original2", "baseURL"),
		entity.NewURL("original3", "baseURL"),
	}
	list := make([]*pb.AddrFullInfoListItemResponse, 0, len(urlList))
	for _, item := range urlList {
		list = append(list, &pb.AddrFullInfoListItemResponse{
			ShortAddr:    item.Short,
			OriginalAddr: item.Original,
		})
	}

	out := &pb.AddrFullInfoListResponse{
		List: list,
	}

	repoPkg.TmpUserID = repoPkg.RightUser
	repoPkg.TmpURLList = urlList

	result, err := h.GetUserURLList(ctx, in)
	assert.Nil(t, err)
	assert.Equal(t, out, result)

	in.User = repoPkg.WrongUser
	_, err = h.GetUserURLList(ctx, in)
	assert.NotNil(t, err)

}
