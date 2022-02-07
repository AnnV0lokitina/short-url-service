package usecase

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

type Repo interface {
	GetInfo() *entity.URLCollection
	SetInfo(info *entity.URLCollection)
}

type UsecaseInterface interface {
	SetURL(full string) (string, string)
	GetURL(uuid string) (string, string, error)
}

type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}
