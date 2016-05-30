package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"testing"
)

func InfiniteIncrementer(step int) stream.GeneratorFunc {
	x := 0
	return stream.GeneratorFunc(func() context.Context {
		ctx := NewContext(context.Background(), x)
		x = x + step
		return ctx
	})
}

func TestTake(t *testing.T) {
	infiniteStreamOfInts, cls := stream.Forever(InfiniteIncrementer(1))

	p := Take(10)
	out := p(infiniteStreamOfInts)
	count := 0

	for value := range out.Values() {
		got := FromContext(value)

		if got != count {
			t.Errorf("Want %d, got %d", count, got)
		}

		count++
	}

	cls()

	if count != 10 {
		t.Errorf("Want %d, got %d", 10, count)
	}
}

func TestTakeUntil(t *testing.T) {
	max := 10
	infiniteStreamOfInts, cls := stream.Forever(InfiniteIncrementer(1))

	p := TakeUntil(Predicate(func(ctx context.Context) bool {
		return FromContext(ctx) >= max
	}))

	out := p(infiniteStreamOfInts)
	count := 0

	for value := range out.Values() {
		got := FromContext(value)

		if got != count {
			t.Errorf("Unexpected value. Want %d, got %d", count, got)
		}

		count++
	}

	cls()

	if count != 10 {
		t.Errorf("Want %d, got %d", 10, count)
	}
}

func TestTakeWhile(t *testing.T) {
	max := 10
	infiniteStreamOfInts, cls := stream.Forever(InfiniteIncrementer(1))

	p := TakeWhile(Predicate(func(ctx context.Context) bool {
		return FromContext(ctx) < max
	}))

	out := p(infiniteStreamOfInts)
	count := 0

	for value := range out.Values() {
		got := FromContext(value)

		if got != count {
			t.Errorf("Unexpected value. Want %d, got %d", count, got)
		}

		count++
	}

	cls()

	if count != 10 {
		t.Errorf("Want %d, got %d", 10, count)
	}
}
