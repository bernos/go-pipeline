package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

// Handler implements a step in our pipeline. A handler takes an incoming context, operates on
// it, and then sends its output to the next Handler in pipeline through the output channel. Any
// errors that are encountered should be sent on the error channel
type Handler interface {
	Handle(context.Context, stream.Stream)
}

// StageFunc makes a regular func implement the Stage interface
type HandlerFunc func(context.Context, stream.Stream)

// Handle satisfies the Stage interface for StageFunc
func (fn HandlerFunc) Handle(ctx context.Context, s stream.Stream) {
	fn(ctx, s)
}
