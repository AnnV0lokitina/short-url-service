package file

import (
	"encoding/json"
	"github.com/AnnV0lokitina/short-url-service.git/internal/entity"
	"os"
)

type Writer struct {
	file    *os.File
	encoder *json.Encoder
}

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

func (w *Writer) WriteURL(url *entity.URL) error {
	return w.encoder.Encode(url)
}
func (w *Writer) Close() error {
	return w.file.Close()
}
