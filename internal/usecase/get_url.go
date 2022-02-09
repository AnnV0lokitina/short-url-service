package usecase

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

func (u *Usecase) GetURL(shortURL string) (*entity.URL, error) {
	url, err := u.repo.GetURL(shortURL)
	if err != nil {
		return nil, err
	}
	return url, nil
}
