package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

// Predicate tests whether a context satisfies some test
type Predicate func(context.Context) bool

// Filter filters values from the input channel that satisfy predicate and sends them
// on the output channel
func Filter(p Predicate) Pipeline {
	return func(in stream.Stream) stream.Stream {
		out, cls := in.WithValues(make(chan context.Context))

		go func() {
			defer cls()

			for ctx := range in.Values() {
				if p(ctx) {
					out.Value(ctx)
				}
			}
		}()

		return out
	}
}
