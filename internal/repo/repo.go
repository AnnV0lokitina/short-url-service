package repo

import (
	"context"
	"errors"
	"sort"

	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
)

// PingBD Show that database is available.
func (r *Repo) PingBD(_ context.Context) bool {
	return true
}

// Close Closes file writer if information stored in file.
func (r *Repo) Close(_ context.Context) error {
	if r.writer != nil {
		return r.writer.Close()
	}
	return nil
}

// SetURL Save url information to storage.
func (r *Repo) SetURL(_ context.Context, userID uint32, url *entity.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()
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

// GetURL Get url information from storage.
func (r *Repo) GetURL(_ context.Context, shortURL string) (*entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
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

// GetUserURLList Get list of urls, created by user, from storage
func (r *Repo) GetUserURLList(_ context.Context, id uint32) ([]*entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	keys := make([]string, 0, len(r.rows))
	for k, row := range r.rows {
		if id != row.UserID {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	logLength := len(keys)
	if logLength <= 0 {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	}
	userLog := make([]*entity.URL, 0, logLength)
	for _, k := range keys {
		userLog = append(userLog, &entity.URL{
			Short:    r.rows[k].ShortURL,
			Original: r.rows[k].OriginalURL,
		})
	}
	return userLog, nil
}

// AddBatch Save to storage list of urls.
func (r *Repo) AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range list {
		record := &entity.Record{
			ShortURL:    item.URL.Short,
			OriginalURL: item.URL.Original,
			UserID:      userID,
			Deleted:     false,
		}
		if r.writer != nil {
			if err := r.writer.WriteRecord(record); err != nil {
				return err
			}
		}
		r.rows[item.URL.Short] = record
	}
	return nil
}

// DeleteBatch Mark urls list like deleted.
func (r *Repo) DeleteBatch(_ context.Context, userID uint32, listShortURL []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, shortURL := range listShortURL {
		record, ok := r.rows[shortURL]
		if ok && record.UserID == userID {
			r.rows[shortURL].Deleted = true
		}
	}
	return nil
}

// CheckUserBatch Return only urls witch can be deleted by user.
func (r *Repo) CheckUserBatch(_ context.Context, userID uint32, listShortURL []string) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	resultList := make([]string, 0, len(listShortURL))
	for _, shortURL := range listShortURL {
		record, ok := r.rows[shortURL]
		if ok && record.UserID == userID {
			resultList = append(resultList, shortURL)
		}
	}
	return resultList, nil
}
