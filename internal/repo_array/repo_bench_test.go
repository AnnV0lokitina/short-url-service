package repoarray

import (
	"context"
	"testing"

	"github.com/AnnV0lokitina/short-url-service/internal/entity"
)

const (
	urlFullString = "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	shortURLHost  = "http://localhost:8080"
)

func BenchmarkRepo_GetURL(b *testing.B) {
	type fields struct {
		list []*entity.Record
	}
	type args struct {
		shortURL string
	}

	url := entity.NewURL(urlFullString, shortURLHost)
	deletedURL := entity.NewURL("deleted_full", shortURLHost)
	list := []*entity.Record{
		{
			OriginalURL: url.Original,
			ShortURL:    url.Short,
			UserID:      1234,
			Deleted:     false,
		},
		{
			OriginalURL: deletedURL.Original,
			ShortURL:    deletedURL.Short,
			UserID:      12345,
			Deleted:     true,
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.URL
		found   bool
		errType string
	}{
		{
			name: "benchmark repo get url",
			fields: fields{
				list: list,
			},
			args: args{
				shortURL: url.Short,
			},
		},
		{
			name: "benchmark repo get url error",
			fields: fields{
				list: list,
			},
			args: args{
				shortURL: "invalid url",
			},
		},
		{
			name: "benchmark repo get url error (empty uuid)",
			fields: fields{
				list: list,
			},
			args: args{
				shortURL: "",
			},
		},
		{
			name: "benchmark repo get deleted url",
			fields: fields{
				list: list,
			},
			args: args{
				shortURL: deletedURL.Short,
			},
		},
	}
	for _, bb := range tests {
		b.Run(bb.name, func(b *testing.B) {
			r := &Repo{
				rows: bb.fields.list,
			}
			for i := 0; i < b.N; i++ {
				r.GetURL(context.TODO(), bb.args.shortURL)
			}
		})
	}
}

func BenchmarkRepo_SetURL(b *testing.B) {
	type fields struct {
		rows []*entity.Record
	}
	type args struct {
		url      *entity.URL
		shortURL string
		userID   uint32
	}

	url := entity.NewURL(urlFullString, shortURLHost)

	tests := []struct {
		name   string
		fields fields
		args   args
		found  bool
	}{
		{
			name: "benchmark set url positive",
			fields: fields{
				rows: make([]*entity.Record, 0),
			},
			args: args{
				url:      url,
				shortURL: url.Short,
				userID:   11,
			},
		},
	}
	for _, bb := range tests {
		b.Run(bb.name, func(b *testing.B) {
			r := &Repo{
				rows: bb.fields.rows,
			}
			for i := 0; i < b.N; i++ {
				r.SetURL(context.TODO(), bb.args.userID, bb.args.url)
			}
		})
	}
}

func BenchmarkRepo_GetUserURLList(b *testing.B) {
	type input struct {
		rows   []*entity.Record
		userID uint32
	}
	rows := []*entity.Record{
		{
			ShortURL:    "short",
			OriginalURL: "original",
			UserID:      1234,
			Deleted:     false,
		},
	}
	tests := []struct {
		name  string
		input input
		want  []*entity.URL
		want1 bool
	}{
		{
			name: "benchmark get urls success",
			input: input{
				rows:   rows,
				userID: 1234,
			},
		},
		{
			name: "benchmark get urls not found",
			input: input{
				rows:   rows,
				userID: 12345,
			},
		},
	}
	for _, bb := range tests {
		b.Run(bb.name, func(b *testing.B) {
			r := &Repo{
				rows: bb.input.rows,
			}
			for i := 0; i < b.N; i++ {
				r.GetUserURLList(context.TODO(), bb.input.userID)
			}
		})
	}
}

func BenchmarkRepo_DeleteBatch(b *testing.B) {
	rows := []*entity.Record{
		{
			ShortURL:    "to_delete",
			OriginalURL: "original",
			UserID:      1234,
			Deleted:     false,
		},
		{
			ShortURL:    "not_delete",
			OriginalURL: "original",
			UserID:      12345,
			Deleted:     false,
		},
	}
	repo := &Repo{
		rows: rows,
	}
	list := []string{"to_delete", "not_delete"}
	for i := 0; i < b.N; i++ {
		repo.DeleteBatch(context.TODO(), 1234, list)
	}
}

func BenchmarkRepo_CheckUserBatch(b *testing.B) {
	rows := []*entity.Record{
		{
			ShortURL:    "to_delete",
			OriginalURL: "original",
			UserID:      1234,
			Deleted:     false,
		},
		{
			ShortURL:    "not_delete",
			OriginalURL: "original",
			UserID:      12345,
			Deleted:     false,
		},
	}
	repo := &Repo{
		rows: rows,
	}
	list := []string{"to_delete", "not_delete"}
	for i := 0; i < b.N; i++ {
		repo.CheckUserBatch(context.TODO(), 1234, list)
	}
}
