package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

// Sink creates a Pipeline that sends all input to fn, and swallows its output
func Sink(fn func(ctx context.Context) error) Pipeline {
	return func(in <-chan context.Context) stream.Stream {
		s := stream.NewStream()

		go func() {
			defer s.Close()

			for ctx := range in {
				if err := fn(ctx); err != nil {
					s.Error(err)
				}
			}
		}()

		return s
	}
}
