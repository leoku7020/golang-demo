package envkit

// EnvOptions is an alias for functional argument.
type EnvOptions func(opts *envOptions)

// envOptions contains all options.
type envOptions struct {
	skip int
}

// WithSkip sets up the skip which is the number of stack frames to ascend
func WithSkip(skip int) EnvOptions {
	return func(opts *envOptions) {
		opts.skip = skip
	}
}

func loadEnvOptions(options ...EnvOptions) *envOptions {
	opts := &envOptions{
		// `skip` is the number of stack frames to ascend, with 0 identifying the caller of Caller.
		// Passing 1 will skip your function, and return the function that called your function.
		skip: 1,
	}
	for _, option := range options {
		option(opts)
	}

	return opts
}
