package pipeline

import (
	"golang.org/x/net/context"
)

type Generator func(context.Context, chan context.Context, chan error)

func Generate(generator Generator) Pipeline {
	return func(in <-chan context.Context) (<-chan context.Context, <-chan error) {
		var (
			out    = make(chan context.Context)
			errors = make(chan error)
		)

		go func() {
			defer close(out)
			defer close(errors)

			for ctx := range in {
				generator(ctx, out, errors)
			}
		}()

		return out, errors
	}
}
