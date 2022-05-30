package file

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
)

const testWriterFileName = "/test_writer"

func TestNewWriter(t *testing.T) {
	type args struct {
		filePath string
		record   *entity.Record
	}
	type resultInterface interface {
		WriteRecord(record *entity.Record) error
		Close() error
	}
	type want struct {
		resultType      string
		interfaceObject interface{}
		record          string
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
				record: &entity.Record{
					UserID:      1234,
					Deleted:     false,
					ShortURL:    "server/checksum",
					OriginalURL: "full",
				},
			},
			want: want{
				resultType:      "*file.Writer",
				interfaceObject: (*resultInterface)(nil),
				record: "{\"user_id\":1234,\"deleted\":false,\"short_url\":\"server/checksum\"," +
					"\"original_url\":\"full\"}\n",
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
			tt.wantURLErr(t, w.WriteRecord(tt.args.record), fmt.Sprintf("WriteRecord(%v)", tt.args.record))
			tt.wantCloseErr(t, w.Close(), "Close()")

			file, err := os.Open(tt.args.filePath)
			require.NoError(t, err)
			data := make([]byte, 100)
			count, err := file.Read(data)
			require.NoError(t, err)
			assert.Equalf(t, tt.want.record, string(data[:count]), "WriteURL(%v)", tt.args.record)
			err = file.Close()
			require.NoError(t, err)
			os.Remove(tt.args.filePath)
		})
	}
	os.RemoveAll(testDir)
}
