package pipeline

import (
	"golang.org/x/net/context"
)

// Pipeline connects an input chan to an output chan. A Pipeline func
// will normally wrap a Stage, and take care of managing channel use, leaving
// the Stage free to concentrate on data manipulation using the context
type Pipeline func(in <-chan context.Context) (<-chan context.Context, <-chan error)

// Run the pipeline using ctx as a starting value. If the context has a timeout or
// deadline, the pipeline will be stopped when it is reached
func (p Pipeline) Run(ctx context.Context) (<-chan context.Context, <-chan error) {
	in := make(chan context.Context)
	done := ctx.Done()

	go func() {
		defer close(in)
		in <- ctx
		<-done
	}()

	return p(in)
}

// Compose sends the output of Pipeline p to the input of Pipeline next
func (p Pipeline) Compose(next Pipeline) Pipeline {
	return Compose(next, p)
}

func (p Pipeline) Pipe(next Handler) Pipeline {
	return Compose(Pipe(next), p)
}

func (p Pipeline) Filter(predicate Predicate) Pipeline {
	return Compose(Filter(predicate), p)
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
