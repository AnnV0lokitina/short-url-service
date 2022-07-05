package file

import (
	"encoding/json"
	"os"

	"github.com/AnnV0lokitina/short-url-service/internal/entity"
)

// Writer Store file pointer and decoder to write to file.
type Writer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewWriter Create new writer.
func NewWriter(filePath string) (*Writer, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteRecord Write record to file.
func (w *Writer) WriteRecord(record *entity.Record) error {
	return w.encoder.Encode(record)
}

// Close Stop work with file.
func (w *Writer) Close() error {
	return w.file.Close()
}
