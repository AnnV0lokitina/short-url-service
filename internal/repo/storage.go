package repo

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo/file"
)

type Repo struct {
	rows   map[string]*entity.Record
	writer *file.Writer
}

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
