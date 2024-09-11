package httpkit

import (
	"net/http"
	"net/textproto"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	PathPattern = "path-pattern"
	platform    = "platform"
)

var (
	retPathPattern = renderHeaderKey(gwruntime.MetadataHeaderPrefix + PathPattern)
	reqPlatform    = renderHeaderKey(platform)
)

func renderHeaderKey(key string) string {
	return textproto.CanonicalMIMEHeaderKey(key)
}

func extractPathPattern(w http.ResponseWriter, r *http.Request) string {
	if path := w.Header().Get(retPathPattern); path != "" {
		return path
	}

	return r.URL.Path
}

func extractPlatform(w http.ResponseWriter, r *http.Request) string {
	return r.Header.Get(reqPlatform)
}
