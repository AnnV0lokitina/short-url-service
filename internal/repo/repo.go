package repo

import (
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo/file"
)

type Repo struct {
	list   map[string]string
	writer *file.Writer
}

func createFilledList(filePath string) (*map[string]string, error) {
	reader, err := file.NewReader(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	list := make(map[string]string)
	for reader.HasNext() {
		url, err := reader.ReadURL()
		if err != nil {
			continue
		}
		list[url.GetChecksum()] = url.GetFullURL()
	}
	return &list, nil
}

func NewRepo(filePath *string) (*Repo, error) {
	if filePath == nil {
		return &Repo{
			list:   make(map[string]string),
			writer: nil,
		}, nil
	}
	list, err := createFilledList(*filePath)
	if err != nil {
		return nil, err
	}
	writer, err := file.NewWriter(*filePath)
	if err != nil {
		return nil, err
	}

	return &Repo{
		list:   *list,
		writer: writer,
	}, nil
}

func (r *Repo) Close() error {
	if r.writer != nil {
		return r.writer.Close()
	}
	return nil
}

func (r *Repo) SetURL(url *entity.URL) error {
	if r.writer != nil {
		if err := r.writer.WriteURL(url); err != nil {
			return err
		}
	}
	r.list[url.GetChecksum()] = url.GetFullURL()
	return nil
}

func (r *Repo) GetURL(checksum string) (*entity.URL, error) {
	fullURL, ok := r.list[checksum]
	if !ok {
		return nil, errors.New("no url saved")
	}
	url := entity.NewURL(fullURL, checksum)
	return url, nil
}
