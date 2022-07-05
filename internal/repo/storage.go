package repo

import (
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	"github.com/AnnV0lokitina/short-url-service/internal/repo/file"
)

// Repo store in - memory storage and file - writer.
type Repo struct {
	rows   map[string]*entity.Record
	writer *file.Writer
}

// NewFileRepo create repository to store information in file.
func NewFileRepo(filePath string) (*Repo, error) {
	records, err := createRecords(filePath)
	if err != nil {
		return nil, err
	}
	writer, err := file.NewWriter(filePath)
	if err != nil {
		return nil, err
	}

	return &Repo{
		rows:   *records,
		writer: writer,
	}, nil
}

// NewMemoryRepo create repository to store information in memory.
func NewMemoryRepo() *Repo {
	return &Repo{
		rows:   make(map[string]*entity.Record),
		writer: nil,
	}
}

func createRecords(filePath string) (*map[string]*entity.Record, error) {
	reader, err := file.NewReader(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	list := make(map[string]*entity.Record)
	for reader.HasNext() {
		record, err := reader.ReadRecord()
		if err != nil {
			continue
		}
		list[record.ShortURL] = record
	}
	return &list, nil
}
