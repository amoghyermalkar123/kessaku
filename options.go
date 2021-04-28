package kessaku

type Options struct {
	WithContext bool
	PoolSize    int
	WithBatch   bool
	TaskCache   chan func()
}

type OptionSetter func(o *Options)

func loadOptions(opts ...OptionSetter) *Options {
	newOptionsInstance := new(Options)

	for _, fn := range opts {
		fn(newOptionsInstance)
	}

	return newOptionsInstance
}

func WithContext() OptionSetter {
	return func(o *Options) {
		o.WithContext = true
	}
}

func WithPoolSize(val int) OptionSetter {
	return func(o *Options) {
		o.PoolSize = val
	}
}

func WithBatch() OptionSetter {
	return func(o *Options) {
		o.WithBatch = true
		o.TaskCache = make(chan func())
	}
}
