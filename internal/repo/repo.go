package repo

import (
	"sync"

	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

type Repo struct {
	info *entity.URLCollection
	mu   sync.Mutex
}

func NewRepo() *Repo {
	collection := entity.NewURLCollection()
	return &Repo{
		info: collection,
	}
}

func (r *Repo) GetInfo() *entity.URLCollection {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.info
}

func (r *Repo) SetInfo(info *entity.URLCollection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.info = info
}
