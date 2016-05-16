package pipeline

import (
	"golang.org/x/net/context"
)

// Sink creates a Pipeline that sends all input to fn, and swallows its output
func Sink(fn func(ctx context.Context) error) Pipeline {
	return func(in <-chan context.Context) (<-chan context.Context, <-chan error) {
		var (
			errors = make(chan error)
		)

		go func() {
			defer close(errors)

			for ctx := range in {
				if err := fn(ctx); err != nil {
					errors <- err
				}
			}
		}()

		return nil, errors
	}
}
