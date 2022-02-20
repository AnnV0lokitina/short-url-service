package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	urlFullString = "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	shortURLHost  = "http://localhost:8080"
)

func createMD5Hash(url string) string {
	byteURL := []byte(url)
	idByte := md5.Sum(byteURL)
	return fmt.Sprintf("%x", idByte)
}

func TestNewURL(t *testing.T) {
	type args struct {
		fullURL  string
		checksum string
	}

	tests := []struct {
		name string
		args args
		want *URL
	}{
		{
			name: "test url created",
			args: args{
				fullURL:  urlFullString,
				checksum: createMD5Hash(urlFullString),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewURL(tt.args.fullURL, tt.args.checksum)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}

func TestURL_GetFullURL(t *testing.T) {
	type fields struct {
		checksum string
		full     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get full url",
			fields: fields{
				checksum: "",
				full:     urlFullString,
			},
			want: urlFullString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				checksum: tt.fields.checksum,
				full:     tt.fields.full,
			}
			got := u.GetFullURL()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestURL_GetShortURL(t *testing.T) {
	type fields struct {
		checksum string
		full     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get short url",
			fields: fields{
				checksum: createMD5Hash(urlFullString),
				full:     urlFullString,
			},
			want: shortURLHost + "/" + createMD5Hash(urlFullString),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				checksum: tt.fields.checksum,
				full:     tt.fields.full,
			}
			got := u.GetShortURL(shortURLHost)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestURL_GetChecksum(t *testing.T) {
	type fields struct {
		checksum string
		full     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get checksum",
			fields: fields{
				checksum: createMD5Hash(urlFullString),
				full:     urlFullString,
			},
			want: createMD5Hash(urlFullString),
		},
		{
			name: "get checksum empty",
			fields: fields{
				checksum: "",
				full:     urlFullString,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				checksum: tt.fields.checksum,
				full:     tt.fields.full,
			}
			got := u.GetChecksum()
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_createChecksum(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "create uuid function positive",
			args: args{
				url: urlFullString,
			},
			want: createMD5Hash(urlFullString),
		},
		{
			name: "create uuid function some string",
			args: args{
				url: "1234",
			},
			want: createMD5Hash("1234"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createChecksum(tt.args.url)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestNewURLFromFullLink(t *testing.T) {
	type args struct {
		fullURL string
	}
	tests := []struct {
		name string
		args args
		want *URL
	}{
		{
			name: "create checksum",
			args: args{
				fullURL: urlFullString,
			},
			want: NewURL(urlFullString, createMD5Hash(urlFullString)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := NewURLFromFullLink(tt.args.fullURL)
			assert.Equalf(t, tt.want, url, "NewURLFromFullLink(%v)", tt.args.fullURL)
		})
	}
}
