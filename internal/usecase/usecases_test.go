package usecase

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/AnnV0lokitina/short-url-service.git/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewUsecase(t *testing.T) {
	type args struct {
		repo Repo
	}
	tests := []struct {
		name string
		args args
		want *Usecase
	}{
		{
			name: "test new usecase positive",
			args: args{
				repo: repo.NewRepo(),
			},
			want: &Usecase{
				repo: repo.NewRepo(),
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
		repo Repo
	}
	type args struct {
		uuid string
	}

	fullURL := "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	url := entity.NewURL(fullURL, "")
	url.CreateUUID()

	repoWithURL := repo.NewRepo()
	repoWithURL.SetURL(url)

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
				uuid: url.GetUUID(),
			},
			want:    url,
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
			want:    url,
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
			want:    url,
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

	fullURL := "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	url := entity.NewURL(fullURL, "")
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
				fullURL: url.GetFullURL(),
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
