package httpkit

import (
	"net/http"

	"github.com/airbrake/gobrake/v5"
)

// A BrakeHandler is an HTTP middleware that provides integration with
// Airbrake.
type BrakeHandler struct {
	Notifier *gobrake.Notifier
}

type airbrakeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewHTTPBrake returns a new BrakeHandler for Airbrake integration. Use the Handler and HandleFunc methods to wrap
// existing HTTP handlers.
func NewHTTPBrake(notifier *gobrake.Notifier) *BrakeHandler {
	h := BrakeHandler{notifier}
	return &h
}

// Handle works as a middleware that wraps an existing http.Handler and sends route performance stats
func (h *BrakeHandler) Handle(handler http.Handler) http.Handler {
	return h.handle(handler)
}

// HandleFunc is like Handler, but with a handler function parameter for cases
// where that is convenient. In particular, use it to wrap a handler function
// literal.
//
//	http.Handle(pattern, h.HandleFunc(func (w http.ResponseWriter, r *http.Request) {
//	    // handler code here
//	}))
func (h *BrakeHandler) HandleFunc(handler http.HandlerFunc) http.HandlerFunc {
	return h.handle(handler)
}

func (h *BrakeHandler) handle(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.Notifier == nil {
			handler.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx, routeMetric := gobrake.NewRouteMetric(ctx, r.Method, r.URL.Path) // Starts the timing
		arw := newAirbrakeResponseWriter(w)

		handler.ServeHTTP(arw, r)

		routeMetric.StatusCode = arw.statusCode
		_ = h.Notifier.Routes.Notify(ctx, routeMetric)
	}
}

func newAirbrakeResponseWriter(w http.ResponseWriter) *airbrakeResponseWriter {
	// Returns 200 OK if WriteHeader isn't called
	return &airbrakeResponseWriter{w, http.StatusOK}
}

func (arw *airbrakeResponseWriter) WriteHeader(code int) {
	arw.statusCode = code
	arw.ResponseWriter.WriteHeader(code)
}
