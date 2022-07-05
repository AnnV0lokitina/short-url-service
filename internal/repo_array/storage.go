package repoarray

import (
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
)

// Repo store in - memory storage and file - writer.
type Repo struct {
	rows []*entity.Record
}

// NewMemoryRepo create repository to store information in memory.
func NewMemoryRepo() *Repo {
	return &Repo{
		rows: make([]*entity.Record, 0),
	}
}
