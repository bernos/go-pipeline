package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

// Pipeline connects an input chan to an output chan. A Pipeline func
// will normally wrap a Stage, and take care of managing channel use, leaving
// the Stage free to concentrate on data manipulation using the context
type Pipeline func(stream.Stream) stream.Stream

// Run the pipeline using ctx as a starting value. If the context has a timeout or
// deadline, the pipeline will be stopped when it is reached
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

// Compose sends the output of Pipeline p to the input of Pipeline next
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

func (p Pipeline) Pipe(next Handler) Pipeline {
	return Compose(Pipe(next), p)
}

func (p Pipeline) Filter(predicate Predicate) Pipeline {
	return Compose(Filter(predicate), p)
}

func Compose(f, g Pipeline) Pipeline {
	return func(in stream.Stream) stream.Stream {
		return f(g(in))
	}
}

// Compose creates a new Pipeline by sending the output of g to the input of f
// func _Compose(f, g Pipeline) Pipeline {
// 	return func(in <-chan context.Context) stream.Stream {
// 		gStream := g(in)
// 		fStream := f(gStream.Values())

// 		go func() {
// 			for err := range gStream.Errors() {
// 				fStream.Error(err)
// 			}
// 		}()

// 		return fStream
// 	}
// }
