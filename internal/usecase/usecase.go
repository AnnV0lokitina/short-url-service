package usecase

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

type Repo interface {
	GetInfo() *entity.UrlCollection
	SetInfo(info *entity.UrlCollection)
}

type UsecaseInterface interface {
	SetUrl(full string) (string, string)
	GetUrl(uuid string) (string, string, error)
}

type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}
