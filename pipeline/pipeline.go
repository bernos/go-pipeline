package pipeline

import (
	"golang.org/x/net/context"
)

// Pipeline connects an input chan to an output chan. A Pipeline func
// will normally wrap a Stage, and take care of managing channel use, leaving
// the Stage free to concentrate on data manipulation using the context
type Pipeline func(in <-chan context.Context) (<-chan context.Context, <-chan error)

// Compose sends the output of Pipeline p to the input of Pipeline next
func (p Pipeline) Compose(next Pipeline) Pipeline {
	return Compose(next, p)
}

// Compose creates a new Pipeline by sending the output of g to the input of f
func Compose(f, g Pipeline) Pipeline {
	return func(in <-chan context.Context) (<-chan context.Context, <-chan error) {
		gValues, gErrors := g(in)
		fValues, fErrors := f(gValues)

		return fValues, MergeErrors(fErrors, gErrors)
	}
}

// MergeErrors sends all errors received on channels f and g on the output channel
func MergeErrors(f, g <-chan error) chan error {
	merged := make(chan error)

	go func() {
		defer close(merged)
		for {
			select {
			case merged <- <-f:
			case merged <- <-g:
			default:
				return
			}
		}
	}()

	return merged
}
