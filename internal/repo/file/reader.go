package file

import (
	"encoding/json"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"os"
)

type Reader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewReader(filePath string) (*Reader, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *Reader) HasNext() bool {
	return r.decoder.More()
}

func (r *Reader) ReadURL() (*entity.URL, error) {
	var url entity.URL
	if err := r.decoder.Decode(&url); err != nil {
		return nil, err
	}
	return &url, nil
}
func (r *Reader) Close() error {
	return r.file.Close()
}
