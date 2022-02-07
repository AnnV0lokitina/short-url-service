package entity

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
)

type UrlCollection struct {
	list map[string]*Url
}

func NewUrlCollection() *UrlCollection {
	list := make(map[string]*Url)
	return &UrlCollection{
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

func (c *UrlCollection) Add(fullUrl string) (string, *Url) {
	url := NewUrl(fullUrl)
	uuid := createUUID()
	c.list[uuid] = url

	return uuid, url
}

func (c *UrlCollection) Get(uuid string) (*Url, error) {
	url, ok := c.list[uuid]
	if !ok {
		return nil, errors.New("No url saved")
	}

	return url, nil
}
