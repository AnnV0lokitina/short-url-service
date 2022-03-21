package handler

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/stretchr/testify/mock"
)

var (
	tmpURL      *entity.URL
	tmpURLError = false
	tmpURLList  []*entity.URL
	tmpUserID   uint32
	pingDB      bool
)

type MockedRepo struct {
	mock.Mock
}

func (r *MockedRepo) SetURL(ctx context.Context, userID uint32, url *entity.URL) error {
	if tmpURLError {
		return errors.New("error")
	}
	tmpURL = url
	tmpURLList = []*entity.URL{tmpURL}
	tmpUserID = userID
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
	return pingDB
}
