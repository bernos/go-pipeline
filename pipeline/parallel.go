package pipeline

import (
	"golang.org/x/net/context"
	"sync"
)

// Parallel runs n instances of a Pipeline in parallel, collecting all output and errors
// onto the output channels
func Parallel(pipeline Pipeline, n int) Pipeline {
	return func(in <-chan context.Context) (<-chan context.Context, <-chan error) {
		var (
			wg     sync.WaitGroup
			out    = make(chan context.Context)
			errors = make(chan error)
		)

		for i := 0; i < n; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				pipelineIn := make(chan context.Context)
				pipelineOut, pipelineErrors := pipeline(pipelineIn)

				defer close(pipelineIn)

				wg.Add(1)
				go func() {
					defer wg.Done()
					for err := range pipelineErrors {
						errors <- err
					}
				}()

				wg.Add(1)
				go func() {
					defer wg.Done()
					for value := range pipelineOut {
						out <- value
					}
				}()

				for ctx := range in {
					pipelineIn <- ctx
				}
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
