package entity

import (
	"fmt"

	"crypto/md5"
)

type URL struct {
	checksum string
	full     string
}

func NewURL(fullURL string, checksum string) *URL {
	return &URL{
		checksum: checksum,
		full:     fullURL,
	}
}

func NewURLFromFullLink(fullURL string) *URL {
	checksum := createChecksum(fullURL)
	return NewURL(fullURL, checksum)
}

func (u *URL) GetFullURL() string {
	return u.full
}

func (u *URL) GetShortURL(serverAddress string) string {
	return serverAddress + "/" + u.checksum
}

func (u *URL) GetChecksum() string {
	return u.checksum
}

func createChecksum(url string) string {
	byteURL := []byte(url)
	idByte := md5.Sum(byteURL)
	return fmt.Sprintf("%x", idByte)
}
