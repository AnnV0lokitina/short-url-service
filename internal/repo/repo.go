package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service.git/pkg/error"
)

func (r *Repo) PingBD(_ context.Context) bool {
	return true
}

func (r *Repo) Close(_ context.Context) error {
	if r.writer != nil {
		return r.writer.Close()
	}
	return nil
}

func (r *Repo) SetURL(_ context.Context, userID uint32, url *entity.URL) error {
	record := &entity.Record{
		ShortURL:    url.Short,
		OriginalURL: url.Original,
		UserID:      userID,
		Deleted:     false,
	}
	if r.writer != nil {
		if err := r.writer.WriteRecord(record); err != nil {
			return err
		}
	}
	r.rows[url.Short] = record
	return nil
}

func (r *Repo) GetURL(_ context.Context, shortURL string) (*entity.URL, error) {
	record, ok := r.rows[shortURL]
	if !ok {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	}
	if record.Deleted {
		return nil, labelError.NewLabelError(labelError.TypeGone, errors.New("URL deleted"))
	}
	url := &entity.URL{
		Short:    record.ShortURL,
		Original: record.OriginalURL,
	}
	return url, nil
}

func (r *Repo) GetUserURLList(_ context.Context, id uint32) ([]*entity.URL, error) {
	userLog := make([]*entity.URL, 0, len(r.rows))
	for _, row := range r.rows {
		if id != row.UserID {
			continue
		}
		userLog = append(userLog, &entity.URL{
			Short:    row.ShortURL,
			Original: row.OriginalURL,
		})
	}
	logLength := len(userLog)
	if logLength <= 0 {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	}
	return userLog[:logLength:logLength], nil
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

func (r *Repo) DeleteBatch(_ context.Context, userID uint32, listShortURL []string) error {
	for _, shortURL := range listShortURL {
		record, ok := r.rows[shortURL]
		if ok && record.UserID == userID {
			r.rows[shortURL].Deleted = true
		}
	}
	return nil
}

func (r *Repo) CheckUserBatch(_ context.Context, userID uint32, listShortURL []string) ([]string, error) {
	resultList := make([]string, 0, len(listShortURL))
	for _, shortURL := range listShortURL {
		record, ok := r.rows[shortURL]
		if ok && record.UserID == userID {
			resultList = append(resultList, shortURL)
		}
	}
	return resultList, nil
}
