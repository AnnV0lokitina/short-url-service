package entity

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
)

type URLCollection struct {
	list map[string]*URL
}

func NewURLCollection() *URLCollection {
	list := make(map[string]*URL)
	return &URLCollection{
		list: list,
	}
}

func createUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func (c *URLCollection) Add(fullURL string) (string, *URL) {
	url := NewURL(fullURL)
	uuid := createUUID()
	c.list[uuid] = url

	return uuid, url
}

func (c *URLCollection) Get(uuid string) (*URL, error) {
	url, ok := c.list[uuid]
	if !ok {
		return nil, errors.New("no url saved")
	}

	return url, nil
}
