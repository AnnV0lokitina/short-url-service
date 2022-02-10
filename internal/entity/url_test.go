package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	urlFullString = "http://xfrpm.ru/ovxnqqxiluncj/lqhza6knc6t2m"
	shortURLHost  = "http://localhost:8080/"
)

func createMD5Hash(url string) string {
	byteURL := []byte(url)
	idByte := md5.Sum(byteURL)
	return fmt.Sprintf("%x", idByte)
}

func TestNewURL(t *testing.T) {
	type args struct {
		fullURL string
		uuid    string
	}

	tests := []struct {
		name string
		args args
		want *URL
	}{
		{
			name: "test url created",
			args: args{
				fullURL: urlFullString,
				uuid:    createMD5Hash(urlFullString),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewURL(tt.args.fullURL, tt.args.uuid)
			assert.ObjectsAreEqual(got, tt.want)
		})
	}
}

func TestURL_CreateUUID(t *testing.T) {
	type fields struct {
		uuid string
		full string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "create uuid",
			fields: fields{
				uuid: "",
				full: urlFullString,
			},
			want: createMD5Hash(urlFullString),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				uuid: tt.fields.uuid,
				full: tt.fields.full,
			}
			u.CreateUUID()
			assert.Equal(t, u.GetUUID(), tt.want)
		})
	}
}

func TestURL_GetFullURL(t *testing.T) {
	type fields struct {
		uuid string
		full string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get full url",
			fields: fields{
				uuid: "",
				full: urlFullString,
			},
			want: urlFullString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				uuid: tt.fields.uuid,
				full: tt.fields.full,
			}
			got := u.GetFullURL()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestURL_GetShortURL(t *testing.T) {
	type fields struct {
		uuid string
		full string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get short url",
			fields: fields{
				uuid: createMD5Hash(urlFullString),
				full: urlFullString,
			},
			want: shortURLHost + createMD5Hash(urlFullString),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				uuid: tt.fields.uuid,
				full: tt.fields.full,
			}
			got := u.GetShortURL()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestURL_GetUUID(t *testing.T) {
	type fields struct {
		uuid string
		full string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get uuid",
			fields: fields{
				uuid: createMD5Hash(urlFullString),
				full: urlFullString,
			},
			want: createMD5Hash(urlFullString),
		},
		{
			name: "get uuid empty",
			fields: fields{
				uuid: "",
				full: urlFullString,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				uuid: tt.fields.uuid,
				full: tt.fields.full,
			}
			got := u.GetUUID()
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_createUUID(t *testing.T) {
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
			got := createUUID(tt.args.url)
			assert.Equal(t, got, tt.want)
		})
	}
}
