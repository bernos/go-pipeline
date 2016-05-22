package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"sync"
	"time"
)

// Pipeline consumes values from an input stream, processes them, and then sends
// new values on an output stream. New pipelines can be created by composing
// existing Pipelines
type Pipeline func(stream.Stream) stream.Stream

// Run the pipeline using ctx as a starting value. The pipeline will be stopped when
// ctx is cancelled
func (p Pipeline) Run(ctx context.Context) stream.Stream {
	in := stream.New()
	done := ctx.Done()

	go func() {
		defer in.Close()
		in.Value(ctx)
		<-done
	}()

	return p(in)
}

// Loop will run a pipeline iteratively, in its own go routine, feeding all values received
// from the output stream back to the input stream. The loop wil continue until ctx is cancelled.
// All values and errors received from the output stream will also be echoed to the Stream
// returned by Loop
func (p Pipeline) Loop(ctx context.Context) stream.Stream {
	var (
		wg sync.WaitGroup

		// Expose the raw value channel for our feedback loop, so that we
		// can use it in a select statement
		buf = make(chan context.Context)
		in  = stream.New().WithValues(buf)

		echo = stream.New()
		done = ctx.Done()
		out  = p(in)
	)

	go func() {
		defer in.Close()
		defer echo.Close()

		for {
			select {
			case <-done:
				wg.Wait()
				return
			case ctx := <-out.Values():
				wg.Add(1)
				go func(ctx context.Context) {
					defer wg.Done()
					ticker := time.NewTicker(time.Nanosecond)

					for {
						select {
						case buf <- ctx:
							echo.Value(ctx)
							return
						case <-ticker.C:
						case <-done:
							return
						}
					}
				}(ctx)
			}
		}
	}()

	go func() {
		for err := range out.Errors() {
			echo.Error(err)
		}
	}()

	in.Value(ctx)

	return echo
}

// Compose creates a new pipeline by passing the output of one pipeline to the input
// of the next
func (p Pipeline) Compose(next Pipeline) Pipeline {
	return Compose(next, p)
}

func (p Pipeline) Map(m Mapper) Pipeline {
	return Compose(Map(m), p)
}

func (p Pipeline) PMap(m Mapper, n int) Pipeline {
	return Compose(PMap(m, n), p)
}

func (p Pipeline) FlatMap(m FlatMapper) Pipeline {
	return Compose(FlatMap(m), p)
}

func (p Pipeline) PFlatMap(m FlatMapper, n int) Pipeline {
	return Compose(PFlatMap(m, n), p)
}

func (p Pipeline) Filter(predicate Predicate) Pipeline {
	return Compose(Filter(predicate), p)
}

func (p Pipeline) Take(n int) Pipeline {
	return Compose(Take(n), p)
}

func (p Pipeline) TakeUntil(predicate Predicate) Pipeline {
	return Compose(TakeUntil(predicate), p)
}

func (p Pipeline) TakeWhile(predicate Predicate) Pipeline {
	return Compose(TakeWhile(predicate), p)
}

// Compose creates a new pipeline by passing the output of one pipeline to the input
// of the next
func Compose(f, g Pipeline) Pipeline {
	return func(in stream.Stream) stream.Stream {
		return f(g(in))
	}
}
