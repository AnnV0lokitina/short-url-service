package entity

import (
	"crypto/md5"
	"fmt"
)

type URL struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

func NewURL(originalURL string, serverAddress string) *URL {
	checksum := createChecksum(originalURL)
	return &URL{
		Short:    CreateShortURL(checksum, serverAddress),
		Original: originalURL,
	}
}

func CreateShortURL(checksum string, serverAddress string) string {
	return serverAddress + "/" + checksum
}

func createChecksum(url string) string {
	byteURL := []byte(url)
	idByte := md5.Sum(byteURL)
	return fmt.Sprintf("%x", idByte)
}
