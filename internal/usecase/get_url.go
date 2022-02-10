package usecase

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

func (u *Usecase) GetURL(uuid string) (*entity.URL, error) {
	url, err := u.repo.GetURL(uuid)
	if err != nil {
		return nil, err
	}
	return url, nil
}
