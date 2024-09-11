package httpkit

import (
	"context"
	"net/http"
)

// Sender is interacting with http request
type Sender interface {
	Send(ctx context.Context, method, url string, header http.Header, value interface{}) (statCode int, body []byte, err error)
	SendReturnCSV(ctx context.Context, method, url string, header http.Header, value interface{}) (statCode int, body [][]string, err error)
}
