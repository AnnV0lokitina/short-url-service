package mockedrepo

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
)

var (
	TmpURL      *entity.URL
	TmpURLError = false
	TmpURLList  []*entity.URL
	TmpUserID   uint32
	PingDB      bool
)

type MockedRepo struct {
	mock.Mock
}

const (
	RightUser     = uint32(1)
	WrongUser     = uint32(2)
	UserWithError = uint32(3)
)

func (r *MockedRepo) SetURL(ctx context.Context, userID uint32, url *entity.URL) error {
	if TmpURLError {
		return errors.New("error")
	}
	TmpURL = url
	TmpURLList = []*entity.URL{TmpURL}
	TmpUserID = userID
	if url.Original == "conflict" {
		return labelError.NewLabelError(labelError.TypeConflict, errors.New("URL exists"))
	}
	if url.Original == "internal error" {
		return errors.New("error")
	}
	return nil
}

func (r *MockedRepo) GetURL(ctx context.Context, shortURL string) (*entity.URL, error) {
	if shortURL == TmpURL.Short {
		return TmpURL, nil
	}
	if shortURL == "URL deleted" {
		return nil, labelError.NewLabelError(labelError.TypeGone, errors.New("URL deleted"))
	}
	if shortURL == "Invalid request" {
		return nil, errors.New("Invalid request")
	}
	return nil, errors.New("no url saved")
}

func (r *MockedRepo) GetUserURLList(_ context.Context, id uint32) ([]*entity.URL, error) {
	if TmpUserID == id {
		return TmpURLList, nil
	}
	if id == 1234 {
		return []*entity.URL{
			{
				Short:    "short",
				Original: "original",
			},
		}, nil
	}
	return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
}

func (r *MockedRepo) Close(_ context.Context) error {
	return nil
}

func (r *MockedRepo) PingBD(ctx context.Context) bool {
	return PingDB
}

func (r *MockedRepo) AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error {
	if userID == UserWithError {
		return errors.New("error")
	}
	return nil
}

func (r *MockedRepo) DeleteBatch(_ context.Context, userID uint32, list []string) error {
	if userID == WrongUser {
		return errors.New("error")
	}
	return nil
}

func (r *MockedRepo) CheckUserBatch(ctx context.Context, userID uint32, list []string) ([]string, error) {
	if userID == WrongUser {
		return nil, errors.New("error")
	}
	return list, nil
}

func (r *MockedRepo) GetStats(_ context.Context) (urls int, users int, err error) {
	return 1, 1, nil
}
