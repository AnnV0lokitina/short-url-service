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
		checksum string
	}

	url := entity.NewURLFromFullLink(urlFullString)
	list := make(map[string]string)
	list[url.GetChecksum()] = url.GetFullURL()

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
				checksum: url.GetChecksum(),
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
				checksum: "invalid checksum",
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
				checksum: "",
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
			got, err := r.GetURL(tt.args.checksum)
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
	}
	type args struct {
		url      *entity.URL
		checksum string
	}

	url := entity.NewURLFromFullLink(urlFullString)
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
				url:      url,
				checksum: url.GetChecksum(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				list: tt.fields.list,
			}
			err := r.SetURL(tt.args.url)
			require.NoError(t, err)

			receiveURL, err := r.GetURL(tt.args.checksum)

			if tt.wantErr == true {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.args.url.GetShortURL(shortURLHost), receiveURL.GetShortURL(shortURLHost))
			assert.Equal(t, tt.args.url.GetFullURL(), receiveURL.GetFullURL())
			assert.Equal(t, tt.args.url.GetChecksum(), receiveURL.GetChecksum())
		})
	}
}
