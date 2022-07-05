package file

import (
	"encoding/json"
	"os"

	"github.com/AnnV0lokitina/short-url-service/internal/entity"
)

// Reader Store file pointer and decoder to read file.
type Reader struct {
	file    *os.File
	decoder *json.Decoder
}

// NewReader create Reader.
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

// HasNext Show if file has next string.
func (r *Reader) HasNext() bool {
	return r.decoder.More()
}

// ReadRecord Read record from file.
func (r *Reader) ReadRecord() (*entity.Record, error) {
	var record entity.Record
	if err := r.decoder.Decode(&record); err != nil {
		return nil, err
	}
	return &record, nil
}

// Close Stop work with file.
func (r *Reader) Close() error {
	return r.file.Close()
}
