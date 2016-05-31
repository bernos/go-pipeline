package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

type Reducer interface {
	Reduce(ctx context.Context, accumulator context.Context) (context.Context, error)
}

type ReducerFunc func(ctx context.Context, accumulator context.Context) (context.Context, error)

func (fn ReducerFunc) Reduce(ctx context.Context, accumulator context.Context) (context.Context, error) {
	return fn(ctx, accumulator)
}

func ReduceLeft(r Reducer) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			out, cls    = in.WithValues(make(chan context.Context))
			accumulator context.Context
		)

		go func() {
			defer cls()

			for ctx := range in.Values() {
				if accumulator == nil {
					accumulator = ctx
				} else {
					result, err := r.Reduce(ctx, accumulator)

					if err == nil {
						accumulator = result
					} else {
						out.Error(err)
					}
				}
			}

			out.Value(accumulator)
		}()

		return out
	}
}

func ReduceRight(r Reducer) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			out, cls    = in.WithValues(make(chan context.Context))
			accumulator context.Context
			values      []context.Context
		)

		go func() {
			defer cls()

			for ctx := range in.Values() {
				values = append(values, ctx)
			}

			if len(values) > 0 {
				i := len(values) - 1

				accumulator = values[i]

				i--

				for i >= 0 {
					result, err := r.Reduce(values[i], accumulator)

					if err == nil {
						accumulator = result
					} else {
						out.Error(err)
					}

					i--
				}

				out.Value(accumulator)
			}
		}()

		return out
	}
}
