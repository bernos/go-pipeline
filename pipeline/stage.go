package pipeline

import (
	"golang.org/x/net/context"
)

// Handler implements a step in our pipeline. A handler takes an incoming context, operates on
// it, and then sends its output to the next Handler in pipeline through the output channel. Any
// errors that are encountered should be sent on the error channel
type Handler interface {
	Handle(context.Context, chan<- context.Context, chan<- error)
}

// StageFunc makes a regular func implement the Stage interface
type HandlerFunc func(context.Context, chan<- context.Context, chan<- error)

// Handle satisfies the Stage interface for StageFunc
func (fn HandlerFunc) Handle(ctx context.Context, out chan<- context.Context, errors chan<- error) {
	fn(ctx, out, errors)
}
