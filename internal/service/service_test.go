package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_repo"
)

func TestNewService(t *testing.T) {
	repo := new(repoPkg.MockedRepo)
	r := NewService("baseURL", repo)
	assert.IsType(t, r, &Service{})
	assert.Equal(t, r.GetBaseURL(), "baseURL")
	assert.Equal(t, r.GetRepo(), repo)

	r.SetBaseURL("baseURL1")
	assert.Equal(t, r.GetBaseURL(), "baseURL1")
}

func TestDeleteURLList(t *testing.T) {
	repo := new(repoPkg.MockedRepo)
	s := NewService("baseURL", repo)
	err := s.DeleteURLList(context.TODO(), repoPkg.WrongUser, []string{"123"})
	assert.NotNil(t, err)
	err = s.DeleteURLList(context.TODO(), repoPkg.RightUser, []string{})
	assert.Nil(t, err)

	s.jobChDelete = make(chan *JobDelete, 1)
	err = s.DeleteURLList(context.TODO(), repoPkg.RightUser, []string{"123"})
	assert.Nil(t, err)
	j := <-s.jobChDelete
	assert.IsType(t, j, &JobDelete{})
	assert.Equal(t, j.UserID, repoPkg.RightUser)
	assert.Equal(t, 1, len(j.URLs))
}
