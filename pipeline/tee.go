package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"sync"
)

// Tee splits a pipeline in two. Inputs are sent to the secondary pipeline, as well as forwarded
// on to the next stage in the main pipeline. Forwarding to the secondary pipeline happens in its
// own go routine, so that the main pipeline is not blocked. The output of the secondary pipeline
// is effectively swallowed
func Tee(pipeline Pipeline) Pipeline {
	return func(in <-chan context.Context) stream.Stream {
		var (
			wg         sync.WaitGroup
			s          = stream.NewStream()
			pipelineIn = make(chan context.Context)
		)

		go func() {
			pipeline(pipelineIn)

			defer close(pipelineIn)
			defer s.Close()

			for ctx := range in {
				wg.Add(1)

				go func(ctx context.Context) {
					defer wg.Done()
					pipelineIn <- ctx
				}(ctx)

				s.Value(ctx)
			}

			wg.Wait()
		}()

		return s
	}
}
