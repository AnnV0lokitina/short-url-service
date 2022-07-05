package file

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AnnV0lokitina/short-url-service/internal/entity"
)

const testReaderFileName = "/test_reader"

func TestNewReader(t *testing.T) {
	type args struct {
		filePath    string
		fileContent string
	}
	type resultInterface interface {
		HasNext() bool
		ReadRecord() (*entity.Record, error)
		Close() error
	}
	type want struct {
		resultType      string
		interfaceObject interface{}
		record          *entity.Record
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
			name: "new reader positive",
			args: args{
				filePath: testDir + testReaderFileName,
				fileContent: "{\"user_id\":1234,\"deleted:\":false,\"short_url\":\"server/checksum\"," +
					"\"original_url\":\"full\"}\n",
			},
			want: want{
				resultType:      "*file.Reader",
				interfaceObject: (*resultInterface)(nil),
				record: &entity.Record{
					UserID:      1234,
					Deleted:     false,
					ShortURL:    "server/checksum",
					OriginalURL: "full",
				},
			},
			wantCreateErr: assert.NoError,
			wantURLErr:    assert.NoError,
			wantCloseErr:  assert.NoError,
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

			r, err := NewReader(tt.args.filePath)
			if !tt.wantCreateErr(t, err, fmt.Sprintf("NewReader(%v)", tt.args.filePath)) {
				return
			}
			assert.Equalf(t, tt.want.resultType, reflect.TypeOf(r).String(), "NewReader(%v)", tt.args.filePath)
			assert.Implements(t, tt.want.interfaceObject, r, "Invalid reader interface")
			assert.Equalf(t, true, r.HasNext(), "HasNext()")
			url, err := r.ReadRecord()
			if !tt.wantURLErr(t, err, "ReadURL()") {
				return
			}
			assert.Equalf(t, tt.want.record, url, "ReadURL()")
			assert.Equalf(t, false, r.HasNext(), "HasNext()")
			tt.wantCloseErr(t, r.Close(), "Close()")
			os.Remove(tt.args.filePath)
		})
	}

	os.RemoveAll(testDir)
}
