package pipeline

import (
	"golang.org/x/net/context"
	"sync"
)

// ParallelPipe runs multiple instances of a Stage in parallel, and passes
// values from an input channel to them. The output of each instance of the
// Stage function will be sent to the output channel
func ParallelPipe(stage Stage, parallelism int) Pipeline {
	return func(in <-chan context.Context) (<-chan context.Context, <-chan error) {
		var wg sync.WaitGroup
		wg.Add(parallelism)

		out := make(chan context.Context)
		errors := make(chan error)

		for i := 0; i < parallelism; i++ {
			go func() {
				for ctx := range in {
					if value, err := stage.Handle(ctx); err == nil {
						out <- value
					} else {
						errors <- err
					}
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
