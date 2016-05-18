package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"sync"
)

// Mapper maps an input context to an output context. It is used to map values
// from input stream to an output stream
type Mapper interface {
	Map(context.Context) (context.Context, error)
}

// MapperFunc is a func that implements Mapper
type MapperFunc func(context.Context) (context.Context, error)

// Map satisfies the Mapper interface
func (fn MapperFunc) Map(ctx context.Context) (context.Context, error) {
	return fn(ctx)
}

// Map creates a Pipeline that maps all values from its input stream to its
// output stream via the Mapper m
func Map(m Mapper) Pipeline {
	return PMap(m, 1)
}

// PMap is a parallel implementation of Map. It produces a Pipeline that maps
// all values from its input stream to its output stream via n concurrent instances
// of the Mapper m
func PMap(m Mapper, n int) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			wg  sync.WaitGroup
			out = in.WithValues(make(chan context.Context))
		)

		wg.Add(n)

		for i := 0; i < n; i++ {
			go func() {
				for ctx := range in.Values() {
					value, err := m.Map(ctx)

					if err == nil {
						out.Value(value)
					} else {
						out.Error(err)
					}
				}
				wg.Done()
			}()
		}

		go func() {
			defer out.Close()
			wg.Wait()
		}()

		return out
	}
}
