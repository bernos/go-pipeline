package stream

import (
	"golang.org/x/net/context"
)

// GeneratorFunc is a func that generates contexts
type GeneratorFunc func() context.Context

// Forever creates a stream that uses generator to generate values until the Stream is closed
func Forever(generator GeneratorFunc) (Stream, CloseFunc) {
	var (
		ch                  = make(chan context.Context)
		output, closeOutput = WithValues(ch)
	)

	go func() {
		defer closeOutput()

		for {
			ctx := generator()

			select {
			case <-output.Done():
				return
			default:
				if ctx.Err() == nil {
					ch <- ctx
				}
			}
		}
	}()

	return output, closeOutput
}
