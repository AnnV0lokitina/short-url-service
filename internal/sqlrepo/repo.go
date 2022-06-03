package sqlrepo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service.git/pkg/error"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

var dbPingTimeout = 1 * time.Second

type PgxIface interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type Repo struct {
	conn PgxIface // *pgx.Conn
	//conn *pgx.Conn
}

func NewSQLRepo(ctx context.Context, dsn string) (*Repo, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	sql := "create table IF NOT EXISTS urls (" +
		"user_id bigint," +
		"short_url text not null," +
		"original_url text not null," +
		"deleted boolean not null default false," +
		"unique (original_url)" +
		")"
	if _, err := conn.Exec(ctx, sql); err != nil {
		return nil, err
	}
	return &Repo{
		conn: conn,
	}, nil
}

func (r *Repo) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r *Repo) PingBD(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := r.conn.Ping(ctx); err != nil {
		return false
	}
	return true
}

func (r *Repo) SetURL(ctx context.Context, userID uint32, url *entity.URL) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "INSERT INTO urls (user_id, short_url, original_url)" +
		"VALUES ($1, $2, $3)" +
		"ON CONFLICT (original_url) DO NOTHING"
	result, err := r.conn.Exec(ctx, sql, userID, url.Short, url.Original)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return labelError.NewLabelError(labelError.TypeConflict, errors.New("URL exists"))
	}
	return nil
}

func (r *Repo) GetURL(ctx context.Context, shortURL string) (*entity.URL, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	var originalURL string
	var deleted bool
	sql := "select original_url, deleted from urls where short_url=$1"
	err := r.conn.QueryRow(ctx, sql, shortURL).Scan(&originalURL, &deleted)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
		}
		return nil, errors.New("get url error")
	}
	if deleted {
		return nil, labelError.NewLabelError(labelError.TypeGone, errors.New("URL deleted"))
	}
	url := &entity.URL{
		Short:    shortURL,
		Original: originalURL,
	}
	return url, nil
}

func (r *Repo) GetUserURLList(ctx context.Context, id uint32) ([]*entity.URL, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "select short_url, original_url from urls where user_id=$1"
	rows, _ := r.conn.Query(ctx, sql, id)
	log := make([]*entity.URL, 0)
	for rows.Next() {
		var shortURL string
		var originalURL string
		err := rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, err
		}
		log = append(log, &entity.URL{
			Short:    shortURL,
			Original: originalURL,
		})
	}
	if len(log) == 0 {
		return nil, labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	}
	return log, nil
}

func (r *Repo) AddBatch(ctx context.Context, userID uint32, list []*entity.BatchURLItem) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql := "INSERT INTO urls (user_id, short_url, original_url)" +
		"VALUES ($1, $2, $3)" +
		"ON CONFLICT (original_url) DO NOTHING"

	_, err = tx.Prepare(ctx, "insert", sql)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, item := range list {
		batch.Queue("insert", userID, item.URL.Short, item.URL.Original)
	}

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}
	br.Close()
	tx.Commit(ctx)
	return nil
}

func (r *Repo) DeleteBatch(ctx context.Context, userID uint32, listShortURL []string) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql := "UPDATE urls " +
		"SET deleted=$1 " +
		"WHERE short_url=$2 " +
		"AND user_id=$3"

	_, err = tx.Prepare(ctx, "delete", sql)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, shortURL := range listShortURL {
		batch.Queue("delete", true, shortURL, userID)
	}

	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}
	br.Close()
	tx.Commit(ctx)
	return nil
}

func (r *Repo) CheckUserBatch(ctx context.Context, userID uint32, listShortURL []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	sql := "SELECT short_url " +
		"FROM urls " +
		"WHERE short_url=$1 " +
		"AND deleted=$2 " +
		"AND user_id=$3"

	_, err := r.conn.Prepare(ctx, "query", sql)
	if err != nil {
		return nil, err
	}

	batch := &pgx.Batch{}
	queryCount := len(listShortURL)

	for _, shortURL := range listShortURL {
		batch.Queue("query", shortURL, false, userID)
	}

	br := r.conn.SendBatch(ctx, batch)

	shortURLs := make([]string, 0, len(listShortURL))
	for i := 0; i < queryCount; i++ {
		rows, err := br.Query()
		if err != nil {
			return nil, err
		}

		for k := 0; rows.Next(); k++ {
			var shortURL string
			if err := rows.Scan(&shortURL); err != nil {
				return nil, err
			}
			shortURLs = append(shortURLs, shortURL)
		}

		if rows.Err() != nil {
			return nil, rows.Err()
		}
	}

	br.Close()
	return shortURLs, nil
}
