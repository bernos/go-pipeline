package stream

import (
	"golang.org/x/net/context"
	"testing"
)

func TestForever(t *testing.T) {
	x := 0

	generator := func() context.Context {
		x++
		return NewContext(context.Background(), x)
	}

	values, cls := Forever(generator)

	want := 1

	for ctx := range values.Values() {
		got := FromContext(ctx)

		if got != want {
			t.Errorf("Want %d, got %d", want, got)
			return
		}

		if got == 10 {
			cls()
		}

		want++
	}
}
