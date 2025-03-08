package data

type Option func(opts *Options)

type Options struct {
	registryKey string // registry key
	SpiderId    int64
}

func WithRegistryKey(key string) Option {
	return func(opts *Options) {
		opts.registryKey = key
	}
}
