package mockedservice

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_repo"
	"github.com/AnnV0lokitina/short-url-service/internal/service"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"

	"github.com/stretchr/testify/mock"
)

var mockedBaseURL string

const (
	WrongUser = uint32(2)
)

type MockedService struct {
	mock.Mock
	trustedSubnet string
}

func NewMockedService() *MockedService {
	service := new(MockedService)
	service.trustedSubnet = "101.101.101"
	return service
}

func (s *MockedService) GetBaseURL() string {
	return mockedBaseURL
}

func (s *MockedService) SetBaseURL(baseURL string) {
	mockedBaseURL = baseURL
}

func (s *MockedService) GetRepo() service.Repo {
	repo := new(repoPkg.MockedRepo)
	return repo
}

func (s *MockedService) DeleteURLList(ctx context.Context, userID uint32, _ []string) error {
	if userID == WrongUser {
		return errors.New("error")
	}
	return nil
}

func (s *MockedService) GetStats(ctx context.Context, ipStr string) (entity.Stats, error) {
	if ipStr == "" {
		return entity.Stats{}, labelError.NewLabelError(labelError.TypeForbidden, errors.New("forbidden"))
	}
	return entity.Stats{
		Users: 1,
		URLs:  1,
	}, nil
}
