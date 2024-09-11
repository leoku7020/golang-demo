package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"demo/pkg/logger"
)

// AllowCORS allows Cross Origin Resource Sharing from any origin.
func AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME: we allow all headers for now. In future we should handle this case via infra.
	headers := []string{"*", "authorization", "content-type", "traceid"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	logger.Debug(fmt.Sprintf("preflight request for %s", r.URL.Path))
}
