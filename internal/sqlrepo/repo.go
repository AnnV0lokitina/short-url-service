package sqlrepo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
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
}

type Repo struct {
	conn PgxIface
}

func NewSQLRepo(ctx context.Context, dsn string) (*Repo, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	sql := "create table IF NOT EXISTS urls (" +
		"user_id bigint," +
		"short_url text not null," +
		"original_url text not null," +
		"unique (short_url)" +
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
	if err := r.conn.Ping(ctx); err != nil {
		return false
	}
	return true
}

func (r *Repo) SetURL(ctx context.Context, userID uint32, url *entity.URL) error {
	// DO NOTHING
	sql := "INSERT INTO urls (user_id, short_url, original_url)" +
		"VALUES ($1, $2, $3)" +
		"ON CONFLICT (short_url) DO NOTHING"
	//"returning original_url, user_id"DO UPDATE set original_url=excluded.original_url
	if _, err := r.conn.Exec(ctx, sql, userID, url.Short, url.Original); err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetURL(ctx context.Context, shortURL string) (*entity.URL, bool, error) {
	var originalURL string
	sql := "select original_url from urls where short_url=$1"
	err := r.conn.QueryRow(ctx, sql, shortURL).Scan(&originalURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, errors.New("no url saved")
	}
	url := &entity.URL{
		Short:    shortURL,
		Original: originalURL,
	}
	return url, true, nil
}

func (r *Repo) GetUserURLList(ctx context.Context, id uint32) ([]*entity.URL, bool, error) {
	sql := "select short_url, original_url from urls where user_id=$1"
	rows, _ := r.conn.Query(ctx, sql, id)
	log := make([]*entity.URL, 0)
	for rows.Next() {
		var shortURL string
		var originalURL string
		err := rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, false, err
		}
		log = append(log, &entity.URL{
			Short:    shortURL,
			Original: originalURL,
		})
	}
	if len(log) == 0 {
		return nil, false, nil
	}
	return log, true, nil
}
