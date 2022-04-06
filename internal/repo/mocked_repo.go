package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service.git/pkg/error"
	"github.com/stretchr/testify/mock"
)

var (
	tmpURL      *entity.URL
	TmpURLError = false
	tmpURLList  []*entity.URL
	tmpUserID   uint32
	PingDB      bool
)

type MockedRepo struct {
	mock.Mock
}

func (r *MockedRepo) SetURL(ctx context.Context, userID uint32, url *entity.URL) error {
	if TmpURLError {
		return errors.New("error")
	}
	tmpURL = url
	tmpURLList = []*entity.URL{tmpURL}
	tmpUserID = userID
	if url.Original == "conflict" {
		return labelError.NewLabelError("CONFLICT", errors.New("URL exists"))
	}
	return nil
}

func (r *MockedRepo) GetURL(ctx context.Context, shortURL string) (*entity.URL, bool, error) {
	if shortURL == tmpURL.Short {
		return tmpURL, true, nil
	}
	return nil, false, errors.New("no url saved")
}

func (r *MockedRepo) GetUserURLList(ctx context.Context, id uint32) ([]*entity.URL, bool, error) {
	if tmpUserID == id {
		return tmpURLList, true, nil
	}
	if id == 1234 {
		return []*entity.URL{
			&entity.URL{
				Short:    "short",
				Original: "original",
			},
		}, true, nil
	}
	return nil, false, nil
}

func (r *MockedRepo) Close(_ context.Context) error {
	return nil
}

func (r *MockedRepo) PingBD(ctx context.Context) bool {
	return PingDB
}

func (r *MockedRepo) AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error {
	return nil
}

func (r *MockedRepo) DeleteBatch(_ context.Context, userID uint32, list []string) error {
	return nil
}

func (r *MockedRepo) CheckUserBatch(ctx context.Context, userID uint32, list []string) ([]string, error) {
	return nil, nil
}
