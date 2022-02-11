package usecase

import (
	"errors"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockedRepo struct {
	mock.Mock
}

var tmpURL *entity.URL

func (r *MockedRepo) SetURL(url *entity.URL) {
	tmpURL = url
}

func (r *MockedRepo) GetURL(uuid string) (*entity.URL, error) {
	if uuid == "uuid" {
		fullURL := "fullURL"
		url := entity.NewURL(fullURL, uuid)
		return url, nil
	}
	return nil, errors.New("no url saved")
}

func TestNewUsecase(t *testing.T) {
	type args struct {
		repo *MockedRepo
	}
	repo := new(MockedRepo)
	tests := []struct {
		name string
		args args
		want *Usecase
	}{
		{
			name: "test new usecase positive",
			args: args{
				repo: repo,
			},
			want: &Usecase{
				repo: repo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUsecase(tt.args.repo)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}

func TestUsecase_GetURL(t *testing.T) {
	type fields struct {
		repo *MockedRepo
	}
	type args struct {
		uuid string
	}

	repoWithURL := new(MockedRepo)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.URL
		wantErr bool
	}{
		{
			name: "test get url positive",
			fields: fields{
				repo: repoWithURL,
			},
			args: args{
				uuid: "uuid",
			},
			want:    entity.NewURL("fullURL", "uuid"),
			wantErr: false,
		},
		{
			name: "test get url negative empty uuid",
			fields: fields{
				repo: repoWithURL,
			},
			args: args{
				uuid: "",
			},
			want:    entity.NewURL("fullURL", "uuid"),
			wantErr: true,
		},
		{
			name: "test get url negative illegal uuid",
			fields: fields{
				repo: repoWithURL,
			},
			args: args{
				uuid: "illegal key",
			},
			want:    entity.NewURL("fullURL", "uuid"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				repo: tt.fields.repo,
			}
			got, err := u.GetURL(tt.args.uuid)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}

func TestUsecase_SetURL(t *testing.T) {
	type fields struct {
		repo Repo
	}
	type args struct {
		fullURL string
	}

	url := entity.NewURL("fullURL", "")
	url.CreateUUID()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *entity.URL
	}{
		{
			name: "set url positive",
			fields: fields{
				repo: repo.NewRepo(),
			},
			args: args{
				fullURL: "fullURL",
			},
			want: url,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usecase{
				repo: tt.fields.repo,
			}
			got := u.SetURL(tt.args.fullURL)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}
