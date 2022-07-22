package sqlrepo

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
	args := flag.Args()
	if len(args) == 0 || args[0] != "local" {
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

	// check that url list empty
	_, err = repo.GetUserURLList(ctx, rightUser)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeNotFound, labelErr.Label)

	// check that no url
	_, err = repo.GetURL(ctx, url.Short)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeNotFound, labelErr.Label)

	// check set url success
	err = repo.SetURL(ctx, rightUser, url)
	assert.Nil(t, err)

	// check add batch success
	err = repo.AddBatch(ctx, rightUser, list)
	assert.Nil(t, err)

	// check get one url success
	readURL, err := repo.GetURL(ctx, url.Short)
	assert.Nil(t, err)
	assert.Equal(t, url, readURL)

	// check get url list success
	readURLList, err := repo.GetUserURLList(ctx, rightUser)
	assert.Nil(t, err)
	assert.Equal(t, len(list)+1, len(readURLList))
	for _, item := range readURLList {
		assert.IsType(t, &entity.URL{}, item)
	}

	// check wrong user url list is empry
	_, err = repo.GetUserURLList(ctx, wrongUser)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeNotFound, labelErr.Label)

	// check set user conflict
	err = repo.SetURL(ctx, rightUser, url)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeConflict, labelErr.Label)

	// check get stats
	nURLs, nUsers, err := repo.GetStats(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 5, nURLs)
	assert.Equal(t, 1, nUsers)

	var shortURLList []string
	for _, item := range list {
		shortURLList = append(shortURLList, item.URL.Short)
	}
	shortURLList = append(shortURLList, "wrong short url")

	// check user url batch exists
	readShortURLList, err := repo.CheckUserBatch(ctx, rightUser, shortURLList)
	assert.Nil(t, err)
	assert.Equal(t, len(list), len(readShortURLList))

	// delete batch
	err = repo.DeleteBatch(ctx, rightUser, readShortURLList)
	assert.Nil(t, err)

	// check that url gone
	_, err = repo.GetURL(ctx, list[0].URL.Short)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeGone, labelErr.Label)

	clearDB(ctx, t, repo)
}
