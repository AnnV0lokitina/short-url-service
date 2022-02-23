package file

import (
	"fmt"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
)

const testWriterFileName = "/test_writer"

func TestNewWriter(t *testing.T) {
	type args struct {
		filePath string
		url      *entity.URL
	}
	type resultInterface interface {
		WriteURL(url *entity.URL) error
		Close() error
	}
	type want struct {
		resultType      string
		interfaceObject interface{}
		url             string
	}
	tmpDir := os.TempDir()
	testDir, err := os.MkdirTemp(tmpDir, "test")
	require.NoError(t, err)

	tests := []struct {
		name          string
		args          args
		want          want
		wantCreateErr assert.ErrorAssertionFunc
		wantURLErr    assert.ErrorAssertionFunc
		wantCloseErr  assert.ErrorAssertionFunc
	}{
		{
			name: "create writer",
			args: args{
				filePath: testDir + testWriterFileName,
				url:      entity.NewURL("full", "checksum"),
			},
			want: want{
				resultType:      "*file.Writer",
				interfaceObject: (*resultInterface)(nil),
				url:             "{\"checksum\":\"checksum\",\"full_url\":\"full\"}\n",
			},
			wantCreateErr: assert.NoError,
			wantURLErr:    assert.NoError,
			wantCloseErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewWriter(tt.args.filePath)
			if !tt.wantCreateErr(t, err, fmt.Sprintf("NewWriter(%v)", tt.args.filePath)) {
				return
			}
			assert.Equalf(t, tt.want.resultType, reflect.TypeOf(w).String(), "NewWriter(%v)", tt.args.filePath)
			assert.Implements(t, tt.want.interfaceObject, w, "Invalid writer interface")
			assert.FileExistsf(t, tt.args.filePath, "file path %v", tt.args.filePath)
			tt.wantURLErr(t, w.WriteURL(tt.args.url), fmt.Sprintf("WriteURL(%v)", tt.args.url))
			tt.wantCloseErr(t, w.Close(), "Close()")

			file, err := os.Open(tt.args.filePath)
			require.NoError(t, err)
			data := make([]byte, 100)
			count, err := file.Read(data)
			require.NoError(t, err)
			assert.Equalf(t, tt.want.url, string(data[:count]), "WriteURL(%v)", tt.args.url)
			err = file.Close()
			require.NoError(t, err)
			os.Remove(tt.args.filePath)
		})
	}
	os.RemoveAll(testDir)
}
