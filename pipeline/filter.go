package pipeline

import (
	"golang.org/x/net/context"
)

// Predicate tests whether a context satisfies some test
type Predicate func(context.Context) bool

// Filter filters values from the input channel that satisfy predicate and sends them
// on the output channel
func Filter(p Predicate) Pipeline {
	return func(in <-chan context.Context) (<-chan context.Context, <-chan error) {
		var (
			out    = make(chan context.Context)
			errors = make(chan error)
		)

		go func() {
			defer close(out)
			defer close(errors)

			for ctx := range in {
				if p(ctx) {
					out <- ctx
				}
			}
		}()

		return out, errors
	}
}
