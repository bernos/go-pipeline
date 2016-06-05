package pipeline

import (
	"sync"

	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

func Loop(p Pipeline) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			wg              sync.WaitGroup
			feedback        = make(chan context.Context)
			pipeIn, pipeCls = stream.WithValues(feedback)
			out             = p(pipeIn)
			echo, cls       = stream.New()
			done            = make(chan interface{})
		)

		go func() {
			defer close(done)

			for ctx := range in.Values() {
				pipeIn.Value(ctx)
			}
		}()

		go func() {

			defer func() {
				defer cls()
				defer pipeCls()
				wg.Wait()
			}()

			for {
				select {
				case <-done:
					return
				case ctx := <-out.Values():
					wg.Add(1)

					go func(ctx context.Context) {
						defer wg.Done()

						select {
						case <-done:
							return
						case feedback <- ctx:
							echo.Value(ctx)
							return
						}
					}(ctx)
				}
			}
		}()

		return echo
	}
}
