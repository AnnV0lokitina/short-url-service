package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerContentEncoding) == encoding {
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error()+"_1", http.StatusInternalServerError)
				return
			}
			r.Body = gzr
			defer gzr.Close()
		}
		if !strings.Contains(r.Header.Get(headerAcceptEncoding), encoding) {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, "Gzip error", http.StatusBadRequest)
			return
		}
		defer gz.Close()
		w.Header().Set(headerContentEncoding, encoding)
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
