package entity

import (
	"fmt"

	"crypto/md5"
)

type URL struct {
	uuid string
	full string
}

func NewURL(fullURL string, uuid string) *URL {
	return &URL{
		uuid: uuid,
		full: fullURL,
	}
}

func (u *URL) CreateUUID() {
	u.uuid = createUUID(u.full)
}

func (u *URL) GetFullURL() string {
	return u.full
}

func (u *URL) GetShortURL() string {
	return "http://localhost:8080/" + u.uuid
}

func (u *URL) GetUUID() string {
	return u.uuid
}

func createUUID(url string) string {
	byteURL := []byte(url)
	idByte := md5.Sum(byteURL)
	return fmt.Sprintf("%x", idByte)
}
