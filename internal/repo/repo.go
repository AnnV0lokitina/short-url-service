package repo

import (
	"context"
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
		list:    *list,
		userLog: make(map[uint32][]*entity.URL),
		writer:  writer,
	}, nil
}

func NewMemoryRepo() *Repo {
	return &Repo{
		list:    make(map[string]string),
		userLog: make(map[uint32][]*entity.URL),
		writer:  nil,
	}
}

func (r *Repo) PingBD(_ context.Context) bool {
	return false
}

func (r *Repo) Close(_ context.Context) error {
	if r.writer != nil {
		return r.writer.Close()
	}
	return nil
}

func (r *Repo) SetURL(_ context.Context, userID uint32, url *entity.URL) error {
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

func (r *Repo) GetURL(_ context.Context, shortURL string) (*entity.URL, bool, error) {
	originalURL, ok := r.list[shortURL]
	if !ok {
		return nil, false, nil
	}
	url := &entity.URL{
		Short:    shortURL,
		Original: originalURL,
	}
	return url, true, nil
}

func (r *Repo) GetUserURLList(_ context.Context, id uint32) ([]*entity.URL, bool, error) {
	log, ok := r.userLog[id]
	if !ok {
		return nil, false, nil
	}
	return log, true, nil
}

func (r *Repo) AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error {
	for _, item := range list {
		err := r.SetURL(ctx, userID, item.URL)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) DeleteBatch(ctx context.Context, userID uint32, list []string) error {
	return nil
}

func (r *Repo) CheckUserBatch(ctx context.Context, userID uint32, list []string) ([]string, error) {
	return nil, nil
}
