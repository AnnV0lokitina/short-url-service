// Package entity has all application entities
package entity

import (
	"crypto/md5"
	"fmt"
)

// URL Stores a short and original url pair.
type URL struct {
	Short    string `json:"short_url"`    // short url
	Original string `json:"original_url"` // original url
}

// NewURL Create new URL structure.
func NewURL(originalURL string, serverAddress string) *URL {
	checksum := createChecksum(originalURL)
	return &URL{
		Short:    CreateShortURL(checksum, serverAddress),
		Original: originalURL,
	}
}

// CreateShortURL Create short url from checksum and server address
func CreateShortURL(checksum string, serverAddress string) string {
	return serverAddress + "/" + checksum
}

func createChecksum(url string) string {
	byteURL := []byte(url)
	idByte := md5.Sum(byteURL)
	return fmt.Sprintf("%x", idByte)
}
