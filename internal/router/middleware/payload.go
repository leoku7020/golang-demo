package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"demo/pkg/logger"
)

func PayloadMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r.Body); err != nil {
			http.Error(w, "Failed to copy request body", http.StatusInternalServerError)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(buf.Bytes()))
		r.Header.Set("Custom-Request-Payload", string(buf.Bytes()))

		if buf.Bytes() != nil {
			logger.Ctx(r.Context()).Info(fmt.Sprintf("request body: %s", string(buf.Bytes())))
		}

		h.ServeHTTP(w, r)
	})
}
