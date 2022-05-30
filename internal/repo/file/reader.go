package file

import (
	"encoding/json"
	"os"

	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
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

func (r *Reader) ReadRecord() (*entity.Record, error) {
	var record entity.Record
	if err := r.decoder.Decode(&record); err != nil {
		return nil, err
	}
	return &record, nil
}
func (r *Reader) Close() error {
	return r.file.Close()
}
