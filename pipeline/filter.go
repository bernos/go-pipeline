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
	return func(in <-chan context.Context) stream.Stream {
		s := stream.NewStream()

		go func() {
			defer s.Close()

			for ctx := range in {
				if p(ctx) {
					s.Value(ctx)
				}
			}
		}()

		return s
	}
}
