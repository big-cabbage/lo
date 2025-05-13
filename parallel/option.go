package parallel

import (
	"golang.org/x/sync/errgroup"
)

type options struct {
	Limit int
}

type Option func(opt *options)

// WithLimit control concurrency limit
func WithLimit(p int) Option {
	return func(opt *options) {
		if p > 0 {
			opt.Limit = p
		}
	}
}

func parseOption(opts []Option) options {
	o := options{
		Limit: -1,
	}
	for _, fn := range opts {
		fn(&o)
	}
	return o
}

func getErrGroup(opts []Option) *errgroup.Group {
	o := parseOption(opts)

	eg := &errgroup.Group{}
	eg.SetLimit(o.Limit)

	return eg
}
