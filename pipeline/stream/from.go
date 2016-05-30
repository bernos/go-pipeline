package stream

import (
	"golang.org/x/net/context"
)

// ContextFunc takes a value and wraps it in a Context
type ContextFunc func(interface{}) context.Context

// From creates a stream of values from the provided slice, wrapping each value in
// a Context created by box. The stream will closed once all values have been sent
func From(values []interface{}, box ContextFunc) (Stream, CloseFunc) {
	var (
		ch                  = make(chan context.Context)
		output, closeOutput = WithValues(ch)
	)

	go func() {
		defer closeOutput()

		for i := range values {
			ctx := box(values[i])

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
