package stream

import (
	"golang.org/x/net/context"
)

type Stream interface {
	Value(context.Context)
	Error(error)
	Values() <-chan context.Context
	Errors() <-chan error
	Close()
	WithValues(chan context.Context) Stream
}

type stream struct {
	values chan context.Context
	errors chan error
}

func NewStream() Stream {
	return &stream{
		values: make(chan context.Context),
		errors: make(chan error),
	}
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

func (s *stream) Close() {
	close(s.errors)
	close(s.values)
}

func (s *stream) WithValues(values chan context.Context) Stream {
	// TODO: may need to forward errors, rather than close over them. If an earlier
	// stream closes the error stream it may cause issues when a later stream tries
	// to send an error
	newStream := &stream{
		values: values,
		errors: s.errors,
	}

	return newStream
}
