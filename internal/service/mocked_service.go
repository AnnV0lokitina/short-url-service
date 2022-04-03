package service

import (
	repoPkg "github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/stretchr/testify/mock"
)

var mockedBaseURL string

type MockedService struct {
	mock.Mock
}

func (s *MockedService) GetBaseURL() string {
	return mockedBaseURL
}

func (s *MockedService) SetBaseURL(baseURL string) {
	mockedBaseURL = baseURL
}

func (s *MockedService) GetRepo() Repo {
	repo := new(repoPkg.MockedRepo)
	return repo
}

func (s *MockedService) DeleteURLList(userID uint32, checksums []string) {
}
