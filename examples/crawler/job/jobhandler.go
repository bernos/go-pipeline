package job

import (
	"fmt"
	"golang.org/x/net/context"
)

// Handler is a handler type that interacts directly with our custom Job type, rather than
// needing to interact with a context.Context to retrieve, update or send a Job through the
// pipeline. A JobHandler will receive a Job via `in`, and send one or more updated Job instances
// by calling the `out()` func.
type Handler func(in Job, out func(Job) error) error

// Handle makes our custom JobHandler implement pipeline.Handler
func (h Handler) Handle(ctx context.Context, out chan<- context.Context, errors chan<- error) {
	if j, ok := FromContext(ctx); ok {
		err := h(j, func(j Job) error {
			out <- NewContext(ctx, j)
			return nil
		})

		if err != nil {
			errors <- err
		}
	} else {
		errors <- fmt.Errorf("Unable to find job in context")
	}
}
