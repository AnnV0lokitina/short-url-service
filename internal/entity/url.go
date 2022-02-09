package entity

import (
	"fmt"

	"crypto/md5"
)

type URL struct {
	Short string
	Full  string
}

func NewURL(fullURL string, shortURL string) *URL {
	return &URL{
		Short: shortURL,
		Full:  fullURL,
	}
}

func (u *URL) CreateShortURL() {
	u.Short = createUUID(u.Full)
}

func createUUID(url string) string {
	byteURl := []byte(url)
	idByte := md5.Sum(byteURl)
	return fmt.Sprintf("%x", idByte)
}
