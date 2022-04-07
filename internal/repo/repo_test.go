package repo

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service.git/pkg/error"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const (
	urlFullString      = "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	shortURLHost       = "http://localhost:8080"
	testReaderFileName = "/test_reader"
)

func TestNewMemoryRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer conn.Close(ctx)

	tests := []struct {
		name string
		want *Repo
	}{
		{
			name: "test new repo positive",
			want: &Repo{
				rows:   make(map[string]*entity.Record),
				writer: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMemoryRepo()
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}

func TestNewFileRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer conn.Close(ctx)

	type args struct {
		filePath    string
		fileContent string
	}

	type want struct {
		repo       *Repo
		listLength int
	}

	tmpDir := os.TempDir()
	testDir, err := os.MkdirTemp(tmpDir, "test")
	require.NoError(t, err)

	tests := []struct {
		name string
		args *args
		want want
	}{
		{
			name: "test new repo positive",
			args: &args{
				filePath: testDir + testReaderFileName,
				fileContent: "{\"user_id\":1234,\"deleted:\":false,\"short_url\":\"server/checksum\"," +
					"\"original_url\":\"full\"}\n",
			},
			want: want{
				repo: &Repo{
					rows: map[string]*entity.Record{"short_url": &entity.Record{
						ShortURL:    "server/checksum",
						OriginalURL: "full",
						UserID:      1234,
						Deleted:     false,
					}},
				},
				listLength: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Create(tt.args.filePath)
			require.NoError(t, err)
			_, err = file.Write([]byte(tt.args.fileContent))
			require.NoError(t, err)
			err = file.Close()
			require.NoError(t, err)

			got, err := NewFileRepo(tt.args.filePath)
			require.NoError(t, err)
			assert.ObjectsAreEqual(got, tt.want)
			assert.Equal(t, len(got.rows), 1, "NewRepo(nil)")
			os.Remove(tt.args.filePath)
		})
	}

	os.RemoveAll(testDir)
}

func TestRepo_GetURL(t *testing.T) {
	type fields struct {
		list map[string]*entity.Record
	}
	type args struct {
		shortURL string
	}

	url := entity.NewURL(urlFullString, shortURLHost)
	deletedUrl := entity.NewURL("deleted_full", shortURLHost)
	list := make(map[string]*entity.Record)
	list[url.Short] = &entity.Record{
		OriginalURL: url.Original,
		ShortURL:    url.Short,
		UserID:      1234,
		Deleted:     false,
	}
	list[deletedUrl.Short] = &entity.Record{
		OriginalURL: deletedUrl.Original,
		ShortURL:    deletedUrl.Short,
		UserID:      12345,
		Deleted:     true,
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
			name: "test repo get url",
			fields: fields{
				list: list,
			},
			args: args{
				shortURL: url.Short,
			},
			want:  url,
			found: true,
		},
		{
			name: "test repo get url error",
			fields: fields{
				list: list,
			},
			args: args{
				shortURL: "invalid url",
			},
			want:    url,
			found:   false,
			errType: labelError.TypeNotFound,
		},
		{
			name: "test repo get url error (empty uuid)",
			fields: fields{
				list: list,
			},
			args: args{
				shortURL: "",
			},
			want:    url,
			found:   false,
			errType: labelError.TypeNotFound,
		},
		{
			name: "test repo get deleted url",
			fields: fields{
				list: list,
			},
			args: args{
				shortURL: deletedUrl.Short,
			},
			want:    deletedUrl,
			found:   false,
			errType: labelError.TypeGone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				rows: tt.fields.list,
			}
			got, err := r.GetURL(context.TODO(), tt.args.shortURL)
			if tt.found {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if tt.errType != "" {
					var labelErr *labelError.LabelError
					assert.True(t, errors.As(err, &labelErr))
					assert.Equal(t, labelErr.Label, tt.errType)
				}
			}
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}

func TestRepo_SetURL(t *testing.T) {
	type fields struct {
		rows map[string]*entity.Record
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
			name: "test set url positive",
			fields: fields{
				rows: make(map[string]*entity.Record),
			},
			args: args{
				url:      url,
				shortURL: url.Short,
				userID:   11,
			},
			found: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				rows: tt.fields.rows,
			}
			err := r.SetURL(context.TODO(), tt.args.userID, tt.args.url)
			require.NoError(t, err)

			receiveURL, err := r.GetURL(context.TODO(), tt.args.shortURL)

			if tt.found {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.args.url.Short, receiveURL.Short)
			assert.Equal(t, tt.args.url.Original, receiveURL.Original)
		})
	}
}

func TestRepo_GetUserURLList(t *testing.T) {
	type input struct {
		rows   map[string]*entity.Record
		userID uint32
	}
	rows := make(map[string]*entity.Record)
	rows["short"] = &entity.Record{
		ShortURL:    "short",
		OriginalURL: "original",
		UserID:      1234,
		Deleted:     false,
	}
	tests := []struct {
		name  string
		input input
		want  []*entity.URL
		want1 bool
	}{
		{
			name: "test get urls",
			input: input{
				rows:   rows,
				userID: 1234,
			},
			want: []*entity.URL{
				&entity.URL{
					Short:    "short",
					Original: "original",
				},
			},
			want1: true,
		},
		{
			name: "test get urls",
			input: input{
				rows:   rows,
				userID: 12345,
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				rows: tt.input.rows,
			}
			got, err := r.GetUserURLList(context.TODO(), tt.input.userID)
			assert.Equalf(t, tt.want, got, "GetUserURLList(%v)", tt.input.userID)
			if tt.want1 {
				require.NoError(t, err)
			}
		})
	}
}

func TestRepo_DeleteBatch(t *testing.T) {
	rows := make(map[string]*entity.Record)
	rows["to_delete"] = &entity.Record{
		ShortURL:    "to_delete",
		OriginalURL: "original",
		UserID:      1234,
		Deleted:     false,
	}
	rows["not_delete"] = &entity.Record{
		ShortURL:    "not_delete",
		OriginalURL: "original",
		UserID:      12345,
		Deleted:     false,
	}
	repo := &Repo{
		rows: rows,
	}
	list := []string{"to_delete", "not_delete"}
	repo.DeleteBatch(context.TODO(), 1234, list)
	assert.True(t, rows["to_delete"].Deleted)
	assert.False(t, rows["not_delete"].Deleted)
}

func TestRepo_CheckUserBatch(t *testing.T) {
	rows := make(map[string]*entity.Record)
	rows["to_delete"] = &entity.Record{
		ShortURL:    "to_delete",
		OriginalURL: "original",
		UserID:      1234,
		Deleted:     false,
	}
	rows["not_delete"] = &entity.Record{
		ShortURL:    "not_delete",
		OriginalURL: "original",
		UserID:      12345,
		Deleted:     false,
	}
	repo := &Repo{
		rows: rows,
	}
	list := []string{"to_delete", "not_delete"}
	resultList, _ := repo.CheckUserBatch(context.TODO(), 1234, list)
	assert.Equal(t, 1, len(resultList))
	assert.Equal(t, "to_delete", list[0])
}
