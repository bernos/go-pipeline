package pipeline

import (
	"time"

	"golang.org/x/net/context"

	"github.com/bernos/go-pipeline/pipeline/stream"
)

// Delay creates a pipeline that waits for the specified duration between
// pulling values from its input stream
func Delay(d time.Duration) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			out, cls = in.WithValues(make(chan context.Context))
		)

		go func() {
			defer cls()

			for value := range in.Values() {
				out.Value(value)
				time.Sleep(d)
			}
		}()

		return out
	}
}
