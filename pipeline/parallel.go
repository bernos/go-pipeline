package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"sync"
)

// Parallel runs n instances of a Pipeline in parallel, collecting all output and errors
// onto the output channels
func Parallel(pipeline Pipeline, n int) Pipeline {
	return func(in <-chan context.Context) stream.Stream {
		var (
			wg sync.WaitGroup
			s  = stream.NewStream()
		)

		for i := 0; i < n; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				pipelineIn := make(chan context.Context)
				// pipelineOut, pipelineErrors := pipeline(pipelineIn)
				pipeOut := pipeline(pipelineIn)

				defer close(pipelineIn)

				wg.Add(1)
				go func() {
					defer wg.Done()
					for err := range pipeOut.Errors() {
						s.Error(err)
					}
				}()

				wg.Add(1)
				go func() {
					defer wg.Done()
					for ctx := range pipeOut.Values() {
						s.Value(ctx)
					}
				}()

				for ctx := range in {
					pipelineIn <- ctx
				}
			}()
		}

		go func() {
			defer s.Close()
			wg.Wait()
		}()

		return s
	}
}
