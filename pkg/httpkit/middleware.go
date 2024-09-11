package httpkit

import "net/http"

// Middleware is an entity (customized http handler) that intercepts the server's request/response life cycle.
// In simple words, it is a piece of code that runs before/after the server caters to a request with a response.
type Middleware func(http.Handler) http.Handler

// Chain indicates a slice of middlewares.
type Chain []Middleware

// returns a Slice of middlewares
func NewChain(middlewares ...Middleware) Chain {
	var slice Chain
	return append(slice, middlewares...)
}

func (c Chain) Then(origHandler http.Handler) http.Handler {
	if origHandler == nil {
		origHandler = http.DefaultServeMux
	}

	// Equivalent to m1(m2(m3(originalHandler)))
	// the reference: https://www.alexedwards.net/blog/making-and-using-middleware
	for i := len(c) - 1; i >= 0; i-- {
		origHandler = c[i](origHandler)
	}

	return origHandler
}
