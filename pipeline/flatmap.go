package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"sync"
)

// FlatMapper maps an input context to a slice of output contexts.
type FlatMapper interface {
	FlatMap(context.Context) ([]context.Context, error)
}

// FlatMapperFunc is a func that impelments FlatMapper
type FlatMapperFunc func(context.Context) ([]context.Context, error)

// FlatMap satisfies the FlatMapper interface
func (fn FlatMapperFunc) FlatMap(ctx context.Context) ([]context.Context, error) {
	return fn(ctx)
}

// FlatMap creates a Pipeline that maps all values from its input stream to its
// output stream via the FlatMapper m. Each Context returned by m.FlatMap will be
// sent as a value on the output stream
func FlatMap(m FlatMapper) Pipeline {
	return PFlatMap(m, 1)
}

// PFlatMap creates a Pipeline that maps all values from its input stream to its
// output stream via n concurrent instances of the FlatMapper m. Each Context
// returned by m.FlatMap will be sent as a value on the output stream
func PFlatMap(m FlatMapper, n int) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			wg  sync.WaitGroup
			out = in.WithValues(make(chan context.Context))
		)

		wg.Add(n)

		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				for ctx := range in.Values() {
					values, err := m.FlatMap(ctx)

					if err == nil {
						for v := range values {
							out.Value(values[v])
						}
					} else {
						out.Error(err)
					}
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
