package pipeline

import (
	"golang.org/x/net/context"
	"sync"
)

// Tee splits a pipeline in two. Inputs are sent to the secondary pipeline, as well as forwarded
// on to the next stage in the main pipeline. Forwarding to the secondary pipeline happens in its
// own go routine, so that the main pipeline is not blocked. The output of the secondary pipeline
// is effectively swallowed
func Tee(pipeline Pipeline) Pipeline {
	return func(in <-chan context.Context) (<-chan context.Context, <-chan error) {
		var (
			wg         sync.WaitGroup
			out        = make(chan context.Context)
			errors     = make(chan error)
			pipelineIn = make(chan context.Context)
		)

		go func() {
			pipeline(pipelineIn)

			defer close(pipelineIn)
			defer close(out)
			defer close(errors)

			for ctx := range in {
				wg.Add(1)

				go func(ctx context.Context) {
					defer wg.Done()
					pipelineIn <- ctx
				}(ctx)

				out <- ctx
			}

			wg.Wait()
		}()

		return out, errors
	}
}
