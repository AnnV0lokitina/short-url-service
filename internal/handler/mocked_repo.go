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

func (r *MockedRepo) SetURL(userID uint32, url *entity.URL) error {
	if tmpURLError {
		return errors.New("error")
	}
	tmpURL = url
	tmpURLList = []*entity.URL{tmpURL}
	tmpUserID = userID
	return nil
}

func (r *MockedRepo) GetURL(shortURL string) (*entity.URL, error) {
	if shortURL == tmpURL.Short {
		return tmpURL, nil
	}
	return nil, errors.New("no url saved")
}

func (r *MockedRepo) GetUserURLList(id uint32) ([]*entity.URL, bool) {
	if tmpUserID == id {
		return tmpURLList, true
	}
	if id == 1234 {
		return []*entity.URL{
			&entity.URL{
				Short:    "short",
				Original: "original",
			},
		}, true
	}
	return nil, false
}

func (r *MockedRepo) PingBD(ctx context.Context) bool {
	return pingDB
}
