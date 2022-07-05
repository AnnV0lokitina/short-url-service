package sqlrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	rightUser    = uint32(1)
	wrongUser    = uint32(2)
	shortURLHost = "http://localhost:8080"
)

type fileConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DataBaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

//func TestRepo_PingBD(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//
//	conn, err := pgxmock.NewConn()
//	require.NoError(t, err)
//	defer conn.Close(ctx)
//
//	type input struct {
//		conn PgxIface
//		ctx  context.Context
//	}
//	tests := []struct {
//		name  string
//		input input
//		want  bool
//	}{
//		{
//			name: "ping positive",
//			input: input{
//				ctx:  ctx,
//				conn: conn,
//			},
//			want: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Repo{
//				conn: tt.input.conn,
//			}
//			assert.Equalf(t, tt.want, r.PingBD(tt.input.ctx), "PingBD(%v)", tt.input.ctx)
//		})
//	}
//}

func clearDB(ctx context.Context, t *testing.T, repo *Repo) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sql := "DELETE FROM urls"
	_, err := repo.conn.Exec(ctx, sql)
	assert.Nil(t, err)
}

func generateBatch() []*entity.BatchURLItem {
	var list []*entity.BatchURLItem
	nums := []int{1, 2, 3, 4}
	for _, num := range nums {
		list = append(list, &entity.BatchURLItem{
			CorrelationID: "add",
			URL:           entity.NewURL(fmt.Sprintf("url%v", num), shortURLHost),
		})
	}
	return list
}

func TestRepo(t *testing.T) {
	if os.Getenv("LOCAL") == "" {
		t.Skip("Skipping testing in CI environment")
	}
	var labelErr *labelError.LabelError

	fContent, err := ioutil.ReadFile("../../cmd/shortener/defaults/defaults_run_test.json")
	assert.Nil(t, err)
	var config = fileConfig{}
	err = json.Unmarshal(fContent, &config)
	assert.Nil(t, err)
	ctx := context.TODO()

	repo, err := NewSQLRepo(ctx, config.DataBaseDSN)
	assert.Nil(t, err)
	defer repo.Close(ctx)

	ping := repo.PingBD(ctx)
	assert.True(t, ping)

	clearDB(ctx, t, repo)

	url := entity.NewURL("url", shortURLHost)
	list := generateBatch()

	_, err = repo.GetUserURLList(ctx, rightUser)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeNotFound, labelErr.Label)

	_, err = repo.GetURL(ctx, url.Short)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeNotFound, labelErr.Label)

	err = repo.SetURL(ctx, rightUser, url)
	assert.Nil(t, err)

	err = repo.AddBatch(ctx, rightUser, list)
	assert.Nil(t, err)

	readURL, err := repo.GetURL(ctx, url.Short)
	assert.Nil(t, err)
	assert.Equal(t, url, readURL)

	readURLList, err := repo.GetUserURLList(ctx, rightUser)
	assert.Nil(t, err)
	assert.Equal(t, len(list)+1, len(readURLList))
	for _, item := range readURLList {
		assert.IsType(t, &entity.URL{}, item)
	}

	_, err = repo.GetUserURLList(ctx, wrongUser)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeNotFound, labelErr.Label)

	err = repo.SetURL(ctx, rightUser, url)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeConflict, labelErr.Label)

	var shortURLList []string
	for _, item := range list {
		shortURLList = append(shortURLList, item.URL.Short)
	}
	shortURLList = append(shortURLList, "wrong short url")
	readShortURLList, err := repo.CheckUserBatch(ctx, rightUser, shortURLList)
	assert.Nil(t, err)
	assert.Equal(t, len(list), len(readShortURLList))

	err = repo.DeleteBatch(ctx, rightUser, readShortURLList)
	assert.Nil(t, err)

	_, err = repo.GetURL(ctx, list[0].URL.Short)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeGone, labelErr.Label)
}
