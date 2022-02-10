package repo

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepo(t *testing.T) {
	list := make(map[string]string)

	tests := []struct {
		name string
		want *Repo
	}{
		{
			name: "test new repo positive",
			want: &Repo{
				list: list,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRepo()
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}

func TestRepo_GetURL(t *testing.T) {
	type fields struct {
		list map[string]string
		mu   sync.Mutex
	}
	type args struct {
		uuid string
	}

	fullURL := "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	url := entity.NewURL(fullURL, "")
	url.CreateUUID()
	list := make(map[string]string)
	list[url.GetUUID()] = url.GetFullURL()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.URL
		wantErr bool
	}{
		{
			name: "test repo get url",
			fields: fields{
				list: list,
			},
			args: args{
				uuid: url.GetUUID(),
			},
			want:    url,
			wantErr: false,
		},
		{
			name: "test repo get url error",
			fields: fields{
				list: list,
			},
			args: args{
				uuid: "invalid uuid",
			},
			want:    url,
			wantErr: true,
		},
		{
			name: "test repo get url error (empty uuid)",
			fields: fields{
				list: list,
			},
			args: args{
				uuid: "",
			},
			want:    url,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				list: tt.fields.list,
				mu:   tt.fields.mu,
			}
			got, err := r.GetURL(tt.args.uuid)
			if tt.wantErr == true {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}

func TestRepo_SetURL(t *testing.T) {
	type fields struct {
		list map[string]string
		mu   sync.Mutex
	}
	type args struct {
		url  *entity.URL
		uuid string
	}

	fullURL := "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	url := entity.NewURL(fullURL, "")
	url.CreateUUID()
	list := make(map[string]string)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test set url positive",
			fields: fields{
				list: list,
			},
			args: args{
				url:  url,
				uuid: url.GetUUID(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				list: tt.fields.list,
				mu:   tt.fields.mu,
			}
			r.SetURL(tt.args.url)

			receiveURL, err := r.GetURL(tt.args.uuid)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.args.url.GetShortURL(), receiveURL.GetShortURL())
			assert.Equal(t, tt.args.url.GetFullURL(), receiveURL.GetFullURL())
			assert.Equal(t, tt.args.url.GetUUID(), receiveURL.GetUUID())
		})
	}
}
