package repo

import (
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	urlFullString      = "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	shortURLHost       = "http://localhost:8080"
	testReaderFileName = "/test_reader"
)

func TestNewMemoryRepo(t *testing.T) {
	list := make(map[string]string)

	tests := []struct {
		name string
		want *Repo
	}{
		{
			name: "test new repo positive",
			want: &Repo{
				list:   list,
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
				filePath:    testDir + testReaderFileName,
				fileContent: "{\"checksum\":\"checksum\",\"full_url\":\"full\"}\n",
			},
			want: want{
				repo: &Repo{
					list: map[string]string{"checksum": "full"},
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
			assert.Equal(t, len(got.list), 1, "NewRepo(nil)")
			os.Remove(tt.args.filePath)
		})
	}

	os.RemoveAll(testDir)
}

func TestRepo_GetURL(t *testing.T) {
	type fields struct {
		list map[string]string
	}
	type args struct {
		shortURL string
	}

	url := entity.NewURL(urlFullString, shortURLHost)
	list := make(map[string]string)
	list[url.Short] = url.Original

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
				shortURL: url.Short,
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
				shortURL: "invalid url",
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
				shortURL: "",
			},
			want:    url,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				list: tt.fields.list,
			}
			got, err := r.GetURL(tt.args.shortURL)
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
		list    map[string]string
		userLog map[uint32][]*entity.URL
	}
	type args struct {
		url      *entity.URL
		shortURL string
		userID   uint32
	}

	url := entity.NewURL(urlFullString, shortURLHost)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test set url positive",
			fields: fields{
				list:    make(map[string]string),
				userLog: make(map[uint32][]*entity.URL),
			},
			args: args{
				url:      url,
				shortURL: url.Short,
				userID:   11,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				list:    tt.fields.list,
				userLog: tt.fields.userLog,
			}
			err := r.SetURL(tt.args.userID, tt.args.url)
			require.NoError(t, err)

			receiveURL, err := r.GetURL(tt.args.shortURL)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.args.url.Short, receiveURL.Short)
			assert.Equal(t, tt.args.url.Original, receiveURL.Original)

			assert.Equal(t, len(r.userLog[tt.args.userID]), 1)
			assert.ObjectsAreEqual(tt.args.url, r.userLog[tt.args.userID][0])
		})
	}
}
