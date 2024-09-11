package handler

import (
	"context"

	"demo/pkg/httpkit"
	"demo/pkg/servkit"
)

func injectHTTPStatus(ctx context.Context, code int) {
	// inject http code if it's under the normal range (200-599)
	if code >= 200 && code < 600 {
		servkit.InjectHTTPCode(ctx, code)
	}
}

func injectHTTPRoute(ctx context.Context) {
	md, ok := servkit.GetMetadata(ctx)
	if !ok {
		return
	}

	if r := md.Route(); r != "" {
		servkit.InjectHTTPHeader(ctx, httpkit.PathPattern, r)
	}
}

// RenderResponse renders response page based on the reference of gPRC response struct.
func RenderResponse[V any](ctx context.Context, ret V, options ...HandlerOptions) (V, error) {
	o := loadErrHandlerOptions(options...)
	injectHTTPStatus(ctx, o.code)
	injectHTTPRoute(ctx)

	// inject headers
	for k, v := range o.headers {
		servkit.InjectHTTPHeader(ctx, k, v)
	}

	return ret, nil
}

// AbortWithError renders error page based on the error and the reference of gPRC response struct,
func AbortWithError[Ptr any](ctx context.Context, err error, _ Ptr, options ...HandlerOptions) (Ptr, error) {
	o := loadErrHandlerOptions(options...)
	injectHTTPStatus(ctx, o.code)
	injectHTTPRoute(ctx)

	// inject headers
	for k, v := range o.headers {
		servkit.InjectHTTPHeader(ctx, k, v)
	}

	var null Ptr
	return null, err
}

type HandlerOptions func(opts *handlerOptions)

type handlerOptions struct {
	code    int
	headers map[string]string
}

// WithHttpStatus injects specified http status code.
func WithHttpStatus(code int) HandlerOptions {
	return func(opts *handlerOptions) {
		opts.code = code
	}
}

// WithHeaders injects specified multiple headers into response headers.
func WithHeaders(headers map[string]string) HandlerOptions {
	return func(opts *handlerOptions) {
		opts.headers = headers
	}
}

func loadErrHandlerOptions(options ...HandlerOptions) *handlerOptions {
	opts := &handlerOptions{}
	for _, option := range options {
		option(opts)
	}

	return opts
}
