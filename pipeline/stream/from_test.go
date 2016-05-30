package stream

import (
	"golang.org/x/net/context"
	"testing"
	"time"
)

type key int

const valuekey key = 0

func NewContext(ctx context.Context, x int) context.Context {
	return context.WithValue(ctx, valuekey, x)
}

func FromContext(ctx context.Context) int {
	return ctx.Value(valuekey).(int)
}

func TestFrom(t *testing.T) {
	values := []interface{}{1, 2, 3, 4, 5}

	box := func(value interface{}) context.Context {
		return NewContext(context.Background(), value.(int))
	}

	out, cls := From(values, box)
	defer cls()

	i := 0

	for ctx := range out.Values() {
		want := values[i].(int)
		got := FromContext(ctx)

		if want != got {
			t.Errorf("Want %d, got %d", want, got)
		}

		i++
	}
}

func TestFromWithContextCancellation(t *testing.T) {
	values := []interface{}{1, 2, 3, 4, 5}

	box := func(value interface{}) context.Context {
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
		ctx = NewContext(ctx, value.(int))
		// wait here, so that context will already have timed out by the time we return
		time.Sleep(time.Millisecond * 10)

		return ctx
	}

	out, cls := From(values, box)
	defer cls()

	i := 0

	for _ = range out.Values() {
		i++
	}

	if i != 0 {
		t.Errorf("Want %d, got %d", 0, i)
	}
}

func TestStreamRequestClose(t *testing.T) {
	values := []interface{}{1, 2, 3, 4, 5}

	box := func(value interface{}) context.Context {
		return NewContext(context.Background(), value.(int))
	}

	out, cls := From(values, box)

	results := make([]int, 0)
	i := 0

	for ctx := range out.Values() {
		want := values[i].(int)
		got := FromContext(ctx)

		if want != got {
			t.Errorf("Want %d, got %d", want, got)
		}

		results = append(results, got)

		if i == 3 {
			cls()
		}

		i++
	}

	if len(results) != 4 {
		t.Errorf("Want length %d, got length %d", 4, len(results))
	}
}
