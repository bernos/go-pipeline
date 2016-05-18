package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"sync"
)

// Pipe passes values from an input channel to a Stage for handling, then
// sends the output of the Stage onto the output channel. The Stage func
// will be run in its own go routine.
func Pipe(stage Handler) Pipeline {
	return ParallelPipe(stage, 1)
}

// ParallelPipe runs multiple instances of a Stage in parallel, and passes
// values from an input channel to them. The output of each instance of the
// Stage function will be sent to the output channel
func ParallelPipe(stage Handler, parallelism int) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var wg sync.WaitGroup
		wg.Add(parallelism)

		out := in.WithValues(make(chan context.Context))

		for i := 0; i < parallelism; i++ {
			go func() {
				for ctx := range in.Values() {
					stage.Handle(ctx, out)
				}
				wg.Done()
			}()
		}

		go func() {
			defer out.Close()
			wg.Wait()
		}()

		return out
	}
}
