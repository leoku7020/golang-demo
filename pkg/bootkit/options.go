package bootkit

// BootOptions sets the options.
type BootOptions interface {
	apply(*bootOptions)
}

type bootOptions struct {
	level int
}

// EmptyBootOption does not alter the server configuration. It can be embedded
// in another structure to build custom server options.
type EmptyBootOption struct{}

func (EmptyBootOption) apply(*bootOptions) {}

// funcBootOption wraps a function that modifies bootOptions into an
// implementation of the BootOptions interface.
type funcBootOption struct {
	f func(*bootOptions)
}

func (o *funcBootOption) apply(do *bootOptions) {
	o.f(do)
}

func newFuncBootOption(f func(*bootOptions)) *funcBootOption {
	return &funcBootOption{
		f: f,
	}
}

// WithShutdownLevel specifies the level used during shutdown process.
// Smaller level will go first. Same level will run concurrently.
// Without specifing the level when calling `AddShutdownHandler()`, it would be 0 by default.
// i.e.,
// AddShutdownHandler(func1, WithShutdownLevel(2))
// AddShutdownHandler(func2, WithShutdownLevel(1))
// AddShutdownHandler(func3, WithShutdownLevel(1))
// func2, func3 go first at the same time. After both finished, func1 go behind.
func WithShutdownLevel(level uint) BootOptions {
	return withShutdownLevel(int(level))
}

func withTopPriorityShutdownLevel() BootOptions {
	return withShutdownLevel(topPriorityShutdownLevel)
}

func withShutdownLevel(level int) BootOptions {
	return newFuncBootOption(func(opts *bootOptions) {
		opts.level = level
	})
}

func applyServOptions(options ...BootOptions) *bootOptions {
	opts := &bootOptions{}
	for _, o := range options {
		o.apply(opts)
	}

	return opts
}
