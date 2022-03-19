package repo

import (
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo/file"
)

type Repo struct {
	list    map[string]string
	userLog map[uint32][]*entity.URL
	writer  *file.Writer
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
		list[url.Short] = url.Original
	}
	return &list, nil
}

func NewFileRepo(filePath string) (*Repo, error) {
	list, err := createFilledList(filePath)
	if err != nil {
		return nil, err
	}
	writer, err := file.NewWriter(filePath)
	if err != nil {
		return nil, err
	}

	return &Repo{
		list:   *list,
		writer: writer,
	}, nil
}

func NewMemoryRepo() *Repo {
	return &Repo{
		list:    make(map[string]string),
		userLog: make(map[uint32][]*entity.URL),
		writer:  nil,
	}
}

func (r *Repo) Close() error {
	if r.writer != nil {
		return r.writer.Close()
	}
	return nil
}

func (r *Repo) SetURL(userID uint32, url *entity.URL) error {
	if r.writer != nil {
		if err := r.writer.WriteURL(url); err != nil {
			return err
		}
	}
	_, exists := r.userLog[userID]
	if !exists {
		r.userLog[userID] = make([]*entity.URL, 0)
	}
	r.userLog[userID] = append(r.userLog[userID], url)
	r.list[url.Short] = url.Original
	return nil
}

func (r *Repo) GetURL(shortURL string) (*entity.URL, error) {
	originalURL, ok := r.list[shortURL]
	if !ok {
		return nil, errors.New("no url saved")
	}
	url := &entity.URL{
		Short:    shortURL,
		Original: originalURL,
	}
	return url, nil
}

func (r *Repo) GetUserURLList(id uint32) ([]*entity.URL, error) {
	log, ok := r.userLog[id]
	if !ok {
		return nil, errors.New("no urls saved")
	}
	return log, nil
}
