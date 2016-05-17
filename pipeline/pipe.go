package pipeline

import (
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
	return func(in <-chan context.Context) (<-chan context.Context, <-chan error) {
		var wg sync.WaitGroup
		wg.Add(parallelism)

		out := make(chan context.Context)
		errors := make(chan error)

		for i := 0; i < parallelism; i++ {
			go func() {
				for ctx := range in {
					stage.Handle(ctx, out, errors)
				}
				wg.Done()
			}()
		}

		go func() {
			defer close(out)
			defer close(errors)
			wg.Wait()
		}()

		return out, errors
	}
}
