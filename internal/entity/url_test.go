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
	userID        = "1234"
)

func createMD5Hash(url string) string {
	byteURL := []byte(url)
	idByte := md5.Sum(byteURL)
	return fmt.Sprintf("%x", idByte)
}

func TestNewRecord(t *testing.T) {
	type args struct {
		userID        string
		originalURL   string
		serverAddress string
	}
	tests := []struct {
		name string
		args args
		want *URL
	}{
		{
			name: "test record created",
			args: args{
				userID:        userID,
				originalURL:   urlFullString,
				serverAddress: shortURLHost,
			},
			want: &URL{
				Short:    shortURLHost + "/" + createMD5Hash(urlFullString),
				Original: urlFullString,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				NewURL(tt.args.originalURL, tt.args.serverAddress),
				"NewURL(%v, %v)",
				tt.args.originalURL,
				tt.args.serverAddress,
			)
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
