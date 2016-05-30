package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

// Sink creates a Pipeline that sends all input to fn, and swallows its output
func Sink(fn func(ctx context.Context) error) Pipeline {
	return func(in stream.Stream) stream.Stream {
		out, cls := stream.New()

		go func() {
			defer cls()

			for ctx := range in.Values() {
				if err := fn(ctx); err != nil {
					out.Error(err)
				}
			}
		}()

		return out
	}
}
