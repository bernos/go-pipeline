package stream

import (
	"fmt"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestStreamWithValuesDoesForwardErrors(t *testing.T) {
	var (
		timeout = time.NewTimer(time.Second)
		ch2     = make(chan context.Context)
		s1, _   = New()
		s2, _   = s1.WithValues(ch2)
		want    = "foo"
	)

	// Send an error on s1
	go func() {
		s1.Error(fmt.Errorf(want))
	}()

	// Assert error appears on output of s2. Timeout after
	// a second to avoid a hanging test
	for {
		select {
		case <-timeout.C:
			t.Error("Timed out")
			return
		case err := <-s2.Errors():
			got := err.Error()

			if want != got {
				t.Errorf("Want %s, got %s", want, got)
			}
			return
		}
	}
}

func TestStreamWithValuesDoesNotForwardValues(t *testing.T) {
	var (
		timeout = time.NewTimer(time.Second)
		ch2     = make(chan context.Context)
		s1, _   = New()
		s2, _   = s1.WithValues(ch2)
		value   = context.WithValue(context.Background(), 0, "foo")
	)

	// Send a value on s1
	go func() {
		s1.Value(value)
	}()

	for {
		select {
		case <-timeout.C:
			return
		case <-s2.Values():
			t.Errorf("Unexpted value from s2")
			return
		}
	}
}
