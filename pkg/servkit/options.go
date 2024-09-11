package servkit

import (
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"demo/pkg/httpkit"
)

// ServOptions sets server options.
type ServOptions interface {
	apply(*servOptions)
}

type servOptions struct {
	middlewares   []httpkit.Middleware
	gwServMuxOpts []gwruntime.ServeMuxOption
	grpcDialOpts  []grpc.DialOption
	grpcServOpts  []grpc.ServerOption
}

// EmptyServOption does not alter the server configuration. It can be embedded
// in another structure to build custom server options.
type EmptyServOption struct{}

func (EmptyServOption) apply(*servOptions) {}

// funcServOption wraps a function that modifies servOptions into an
// implementation of the ServOptions interface.
type funcServOption struct {
	f func(*servOptions)
}

func (o *funcServOption) apply(do *servOptions) {
	o.f(do)
}

func newFuncServOption(f func(*servOptions)) *funcServOption {
	return &funcServOption{
		f: f,
	}
}

// WithMiddlewares specifies the middlewares used in grpc gateway.
// A middleware can do the following:
// - Process the request before running business logic (authentication)
// - Modify the request to the next handler function (attaching payload)
// - Modify the response for the client
// - Logging.... and much more
func WithMiddlewares(middlewares ...httpkit.Middleware) ServOptions {
	return newFuncServOption(func(opts *servOptions) {
		opts.middlewares = middlewares
	})
}

// WithGWServMuxOptions specifies the options used in gwruntime.NewServeMux()
func WithGWServMuxOptions(options ...gwruntime.ServeMuxOption) ServOptions {
	return newFuncServOption(func(opts *servOptions) {
		opts.gwServMuxOpts = options
	})
}

// WithGrpcDialOptions specifies the options used in grpc.DialContext()
func WithGrpcDialOptions(options ...grpc.DialOption) ServOptions {
	return newFuncServOption(func(opts *servOptions) {
		opts.grpcDialOpts = options
	})
}

// WithGrpcServOptions specifies the options used in grpc.NewServer()
func WithGrpcServOptions(options ...grpc.ServerOption) ServOptions {
	return newFuncServOption(func(opts *servOptions) {
		opts.grpcServOpts = options
	})
}

func applyServOptions(options ...ServOptions) *servOptions {
	opts := &servOptions{}
	for _, o := range options {
		o.apply(opts)
	}

	return opts
}
