package repo

import (
	"errors"
	"sync"

	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

type Repo struct {
	list map[string]string
	mu   sync.Mutex
}

func NewRepo() *Repo {
	list := make(map[string]string)
	return &Repo{
		list: list,
	}
}

func (r *Repo) SetURL(url *entity.URL) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.list[url.Short] = url.Full
}

func (r *Repo) GetURL(shortURL string) (*entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fullURL, ok := r.list[shortURL]
	if !ok {
		return nil, errors.New("no url saved")
	}
	url := entity.NewURL(fullURL, shortURL)
	return url, nil
}
