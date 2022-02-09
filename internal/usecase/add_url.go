package usecase

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

func (u *Usecase) SetURL(fullURL string) *entity.URL {
	url := entity.NewURL(fullURL, "")
	url.CreateUUID()
	u.repo.SetURL(url)

	return url
}
