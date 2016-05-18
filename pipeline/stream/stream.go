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
	Close()

	// Create a child Stream, inherriting errors from the parent stream, with the provided
	// values channel
	WithValues(chan context.Context) Stream
}

type stream struct {
	values chan context.Context
	errors chan error
}

func NewStream() Stream {
	return &stream{}
	// return &stream{
	// 	values: make(chan context.Context),
	// 	errors: make(chan error),
	// }
}

func (s *stream) Values() <-chan context.Context {
	return s.getValues()
}

func (s *stream) Errors() <-chan error {
	return s.getErrors()
}

func (s *stream) Value(ctx context.Context) {
	s.getValues() <- ctx
}

func (s *stream) Error(err error) {
	s.getErrors() <- err
}

func (s *stream) Close() {
	if s.errors != nil {
		close(s.errors)
	}
	if s.values != nil {
		close(s.values)
	}
}

func (s *stream) WithValues(values chan context.Context) Stream {
	// TODO: may need to forward errors, rather than close over them. If an earlier
	// stream closes the error stream it may cause issues when a later stream tries
	// to send an error
	newStream := &stream{
		values: values,
		errors: s.getErrors(),
	}

	return newStream
}

func (s *stream) getValues() chan context.Context {
	if s.values == nil {
		s.values = make(chan context.Context)
	}
	return s.values
}

func (s *stream) getErrors() chan error {
	if s.errors == nil {
		s.errors = make(chan error)
	}
	return s.errors
}
