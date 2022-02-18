package repo

import (
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

type Repo struct {
	list map[string]string
}

func NewRepo() *Repo {
	list := make(map[string]string)
	return &Repo{
		list: list,
	}
}

func (r *Repo) SetURL(url *entity.URL) {
	r.list[url.GetChecksum()] = url.GetFullURL()
}

func (r *Repo) GetURL(checksum string) (*entity.URL, error) {
	fullURL, ok := r.list[checksum]
	if !ok {
		return nil, errors.New("no url saved")
	}
	url := entity.NewURL(fullURL, checksum)
	return url, nil
}
