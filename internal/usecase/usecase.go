package usecase

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

type Repo interface {
	SetURL(url *entity.URL)
	GetURL(uuid string) (*entity.URL, error)
}

type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}
