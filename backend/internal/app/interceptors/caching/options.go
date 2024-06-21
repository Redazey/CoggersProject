package caching

import "context"

var (
	defaultOptions = &options{
		cacheHandlerFunc: nil,
	}
)

type options struct {
	cacheHandlerFunc CacheHandlerFuncContext
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

func WithCacheHandler(f CacheHandlerFunc) Option {
	return func(o *options) {
		o.cacheHandlerFunc = CacheHandlerFuncContext(func(ctx context.Context, req interface{}) (res interface{}, err error) {
			return f(req)
		})
	}
}

func WithCacheHandlerContext(f CacheHandlerFuncContext) Option {
	return func(o *options) {
		o.cacheHandlerFunc = f
	}
}
