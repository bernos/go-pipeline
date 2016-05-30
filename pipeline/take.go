package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
)

// Take creates a pipeline that returns at most n items from the input stream,
// and ignores all further values
func Take(n int) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			out, cls = in.WithValues(make(chan context.Context))
		)

		go func() {
			defer cls()

			count := 0

			for value := range in.Values() {
				if count < n {
					out.Value(value)
					count++
				} else {
					return
				}
			}
		}()

		return out
	}
}

// TakeUntil creates a pipeline that will forward values from its input stream
// to its output stream up until a value from the input stream satisfies the
// predicate
func TakeUntil(predicate Predicate) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			out, cls = in.WithValues(make(chan context.Context))
		)

		go func() {
			defer cls()

			for value := range in.Values() {
				if predicate(value) {
					return
				} else {
					out.Value(value)
				}
			}
		}()

		return out
	}
}

// TakeWhile creates a pipeline that will forward values from its input stream
// to its output stream until a value that doesnt satify the predicate is found
func TakeWhile(predicate Predicate) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			out, cls = in.WithValues(make(chan context.Context))
		)

		go func() {
			defer cls()

			for value := range in.Values() {
				if predicate(value) {
					out.Value(value)
				} else {
					return
				}
			}
		}()

		return out
	}
}
