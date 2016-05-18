package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"sync"
)

// Parallel runs n instances of a Pipeline in parallel, collecting all output and errors
// onto the output channels
func Parallel(pipeline Pipeline, n int) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			wg  sync.WaitGroup
			out = in.WithValues(make(chan context.Context))
		)

		for i := 0; i < n; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				pipelineIn := stream.New()
				pipeOut := pipeline(pipelineIn)

				defer pipelineIn.Close()

				wg.Add(1)
				go func() {
					defer wg.Done()
					for err := range pipeOut.Errors() {
						out.Error(err)
					}
				}()

				wg.Add(1)
				go func() {
					defer wg.Done()
					for ctx := range pipeOut.Values() {
						out.Value(ctx)
					}
				}()

				for ctx := range in.Values() {
					pipelineIn.Value(ctx)
				}
			}()
		}

		go func() {
			defer out.Close()
			wg.Wait()
		}()

		return out
	}
}
