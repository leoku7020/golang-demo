package httpkit

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"

	"demo/pkg/logger"
)

// ref: https://github.com/newrelic/go-agent/blob/master/v3/newrelic/instrumentation.go

type NRHandler struct {
	app *newrelic.Application
}

// NewNRHttpHandler returns a new NRHandler for New Relic integration. Use the Handler and HandleFunc methods to wrap
// existing HTTP handlers.
func NewNRHttpHandler(app *newrelic.Application) *NRHandler {
	return &NRHandler{app: app}
}

// Handle works as a middleware that wraps an existing http.Handler and sends data to New Relic.
func (h *NRHandler) Handle(handler http.Handler) http.Handler {
	return h.handle(handler)
}

// HandleFunc is like Handler, but with a handler function parameter for cases
// where that is convenient. In particular, use it to wrap a handler function
// literal.
//
//	http.Handle(pattern, h.HandleFunc(func (w http.ResponseWriter, r *http.Request) {
//	    // handler code here
//	}))
func (h *NRHandler) HandleFunc(handler http.HandlerFunc) http.HandlerFunc {
	return h.handle(handler)
}

func (h *NRHandler) handle(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.app == nil {
			handler.ServeHTTP(w, r)
			return
		}

		path := extractPathPattern(w, r)
		name := r.Method + " " + path

		txn := h.app.StartTransaction(name)
		defer txn.End()

		w = txn.SetWebResponse(w)
		txn.SetWebRequestHTTP(r)

		// inject traceID in context for logger
		ctx := r.Context()
		md := txn.GetTraceMetadata()

		traceID := md.TraceID
		spanID := md.SpanID
		if traceID == "" {
			traceID = strings.Replace(uuid.NewString(), "-", "", -1)
		}
		if spanID == "" {
			spanID = strings.Replace(uuid.NewString(), "-", "", -1)
		}

		ctx = logger.ContextWithFields(ctx, logger.Fields{
			"nr.traceID": traceID,
			"nr.spanID":  spanID,
		})

		logger.Ctx(ctx).Info(name, logger.WithField("http.query", r.URL.RawQuery))

		// inject txn in context
		ctx = newrelic.NewContext(ctx, txn)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	}
}
