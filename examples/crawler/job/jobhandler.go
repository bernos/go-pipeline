package job

import (
	"fmt"
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

// Handler is a handler type that interacts directly with our custom Job type, rather than
// needing to interact with a context.Context to retrieve, update or send a Job through the
// pipeline. A JobHandler will receive a Job via `in`, and send one or more updated Job instances
// by calling the `out()` func.
type Handler func(in Job, out func(Job) error) error

// Handle makes our custom JobHandler implement pipeline.Handler
func (h Handler) Handle(ctx context.Context, out stream.Stream) {
	if j, ok := FromContext(ctx); ok {
		err := h(j, func(j Job) error {
			out.Value(NewContext(ctx, j))
			return nil
		})

		if err != nil {
			out.Error(err)
		}
	} else {
		out.Error(fmt.Errorf("Unable to find job in context"))
	}
}

type Mapper func(in Job) (Job, error)

func (m Mapper) Map(ctx context.Context) (context.Context, error) {
	if j, ok := FromContext(ctx); ok {
		value, err := m(j)

		if err == nil {
			return nil, err
		}

		return NewContext(ctx, value), nil
	} else {
		return nil, fmt.Errorf("Unable to find job in context")
	}
}

type FlatMapper func(in Job) ([]Job, error)

func (m FlatMapper) FlatMap(ctx context.Context) ([]context.Context, error) {
	if j, ok := FromContext(ctx); ok {
		jobs, err := m(j)

		if err == nil {
			return nil, err
		}

		values := make([]context.Context, len(jobs))

		for i := range jobs {
			values[i] = NewContext(ctx, jobs[i])
		}

		return values, nil
	} else {
		return nil, fmt.Errorf("Unable to find job in context")
	}
}
