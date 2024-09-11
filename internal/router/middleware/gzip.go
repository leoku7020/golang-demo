package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"demo/pkg/logger"
)

func GZipDecompressor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if encode := r.Header.Get("Content-Encoding"); encode == "gzip" {
			// Handle gzip encoding
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to create gzip reader", http.StatusBadRequest)
				return
			}
			defer reader.Close()

			// Replace the request body with the decompressed data
			decompressedData, err := io.ReadAll(reader)
			if err != nil {
				http.Error(w, "Failed to read decompressed data", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(decompressedData))

			logger.Ctx(r.Context()).Info(fmt.Sprintf("Decompressed request body: %s", string(decompressedData)))
		}

		h.ServeHTTP(w, r)
	})
}
