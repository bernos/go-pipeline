// Package stream implements a stream used to pass values and errors between pipeline
// functions. Values are wrapped in context.Context, and internally channels are used
// for moving errors and values around
package stream

import (
	"golang.org/x/net/context"
)

// Stream is an interface that facilitates sending values and errors through a pipeline
// Internally the stream uses an error channel, and a Context channel for sending values
type Stream interface {
	// Send a value on the stream
	Value(context.Context)

	// Send an error on the stream
	Error(error)

	// Retrieve the read only values channel
	Values() <-chan context.Context

	// Retrieve the read only error channel
	Errors() <-chan error

	// Close the Stream. This will close both the error and value channels
	// Close()

	// Retrieve the done channel
	Done() <-chan struct{}

	// Create a child Stream, inherriting errors from the parent stream, with the provided
	// values channel
	WithValues(chan context.Context) (Stream, CloseFunc)
}

// A CloseFunc closes a stream. If a CloseFunc is called more than once it does nothing
type CloseFunc func()

type stream struct {
	values chan context.Context
	errors chan error
	done   chan struct{}
	closed bool
}

// New creates an initialized Stream
func New() (Stream, CloseFunc) {
	s := &stream{
		values: make(chan context.Context),
		errors: make(chan error),
		done:   make(chan struct{}),
	}
	return s, closeStream(s)
}

// WithValues creates an initialized Stream from an existing value channel
func WithValues(values chan context.Context) (Stream, CloseFunc) {
	s := &stream{
		values: values,
		errors: make(chan error),
		done:   make(chan struct{}),
	}
	return s, closeStream(s)
}

func (s *stream) Values() <-chan context.Context {
	return s.values
}

func (s *stream) Errors() <-chan error {
	return s.errors
}

func (s *stream) Value(ctx context.Context) {
	s.values <- ctx
}

func (s *stream) Error(err error) {
	s.errors <- err
}

func (s *stream) Done() <-chan struct{} {
	return s.done
}

// WithValues creates a new Stream from an existing value channel. Errors from
// Stream s will be forwarded to the new stream
func (s *stream) WithValues(values chan context.Context) (Stream, CloseFunc) {
	newStream := &stream{
		values: values,
		errors: make(chan error),
		done:   make(chan struct{}),
	}

	go func() {
		for err := range s.Errors() {
			newStream.errors <- err
		}
	}()

	return newStream, closeStream(newStream)
}

func closeStream(s *stream) CloseFunc {
	return func() {
		if !s.closed {
			go func() {
				<-s.done
				close(s.errors)
				close(s.values)
			}()

			close(s.done)
			s.closed = true
		}
	}
}
