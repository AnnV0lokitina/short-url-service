package repoarray

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
)

// SetURL Save url information to storage.
func (r *Repo) SetURL(_ context.Context, userID uint32, url *entity.URL) error {
	record := &entity.Record{
		ShortURL:    url.Short,
		OriginalURL: url.Original,
		UserID:      userID,
		Deleted:     false,
	}
	r.rows = append(r.rows, record)
	return nil
}

func findByShortURL(rows []*entity.Record, shortURL string) (*entity.Record, bool) {
	for i := range rows {
		if rows[i].ShortURL == shortURL {
			return rows[i], true
		}
	}
	return nil, false
}

// GetURL Get url information from storage.
func (r *Repo) GetURL(_ context.Context, shortURL string) (*entity.URL, error) {
	record, ok := findByShortURL(r.rows, shortURL)
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

// GetUserURLList Get list of urls, created by user, from storage
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
	return userLog, nil
}

// AddBatch Save to storage list of urls.
func (r *Repo) AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error {
	for _, item := range list {
		err := r.SetURL(ctx, userID, item.URL)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteBatch Mark urls list like deleted.
func (r *Repo) DeleteBatch(_ context.Context, userID uint32, listShortURL []string) error {
	for _, shortURL := range listShortURL {
		record, ok := findByShortURL(r.rows, shortURL)
		if ok && record.UserID == userID {
			record.Deleted = true
		}
	}
	return nil
}

// CheckUserBatch Return only urls witch can be deleted by user.
func (r *Repo) CheckUserBatch(_ context.Context, userID uint32, listShortURL []string) ([]string, error) {
	resultList := make([]string, 0, len(listShortURL))
	for _, shortURL := range listShortURL {
		record, ok := findByShortURL(r.rows, shortURL)
		if ok && record.UserID == userID {
			resultList = append(resultList, shortURL)
		}
	}
	return resultList, nil
}
