package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"log"
)

// Take creates a pipeline that returns at most n items from the input stream,
// and ignores all further values
func Take(n int) Pipeline {
	return func(in stream.Stream) stream.Stream {
		var (
			out = in.WithValues(make(chan context.Context))
		)

		go func() {
			defer out.Close()
			count := 0

			for value := range in.Values() {
				if count < n {
					out.Value(value)
				}
				count++
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
			out = in.WithValues(make(chan context.Context))
			ok  = true
		)

		go func() {
			defer out.Close()

			for value := range in.Values() {
				if ok {
					if predicate(value) {
						ok = false
					} else {
						out.Value(value)
					}
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
			out = in.WithValues(make(chan context.Context))
			ok  = true
		)

		go func() {
			defer out.Close()

			for value := range in.Values() {
				if ok {
					if predicate(value) {
						out.Value(value)
					} else {
						ok = false
					}
				}
			}
		}()

		return out
	}
}
