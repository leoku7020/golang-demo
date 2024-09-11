package logger

// LoggerOptions sets the options.
type LoggerOptions interface {
	apply(*loggerOptions)
}

type pair struct {
	Key string
	Val interface{}
}

type loggerOptions struct {
	pairs  []pair
	fields Fields
	err    error
}

// funcLoggerOption wraps a function that modifies loggerOptions into an
// implementation of the LoggerOptions interface.
type funcLoggerOption struct {
	f func(*loggerOptions)
}

func (o *funcLoggerOption) apply(do *loggerOptions) {
	o.f(do)
}

func newFuncLoggerOption(f func(*loggerOptions)) *funcLoggerOption {
	return &funcLoggerOption{
		f: f,
	}
}

// WithFields specifies the fields (key-value pairs) when logging messages.
func WithFields(fields Fields) LoggerOptions {
	return newFuncLoggerOption(func(opts *loggerOptions) {
		opts.fields = fields
	})
}

// WithFields specifies the field (key-value pairs) when logging messages.
func WithField(key string, val interface{}) LoggerOptions {
	return newFuncLoggerOption(func(opts *loggerOptions) {
		opts.pairs = append(opts.pairs, pair{Key: key, Val: val})
	})
}

// WithError specifies the error with the key `err` when logging messages.
func WithError(err error) LoggerOptions {
	return newFuncLoggerOption(func(opts *loggerOptions) {
		opts.err = err
	})
}

func applyLoggerOptions(options ...LoggerOptions) *loggerOptions {
	opts := &loggerOptions{}
	for _, o := range options {
		o.apply(opts)
	}

	return opts
}
