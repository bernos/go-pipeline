package pipeline

import (
	"fmt"
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"testing"
	"time"
)

type IntMapper func(x int) int

func (m IntMapper) Map(ctx context.Context) (context.Context, error) {
	return NewContext(ctx, m(FromContext(ctx))), nil
}

func DelayedMultiplyMapper(n int, d time.Duration) IntMapper {
	return IntMapper(func(x int) int {
		time.Sleep(d)
		return x * n
	})
}

func TestMap(t *testing.T) {
	a := 2
	b := 2
	pl := Map(DelayedMultiplyMapper(a, time.Millisecond))
	input, cls := stream.New()
	output := pl(input)

	defer cls()

	go func() {
		input.Value(NewContext(context.Background(), b))
	}()

	result := <-output.Values()

	want := a * b
	got := FromContext(result)

	if got != want {
		t.Errorf("Want %d, got %d", want, got)
	}
}

func TestMapError(t *testing.T) {
	mapper := MapperFunc(func(ctx context.Context) (context.Context, error) {
		return nil, fmt.Errorf("An error")
	})

	pl := Map(mapper)
	input, cls := stream.New()
	output := pl(input)

	defer cls()

	go func() {
		input.Value(NewContext(context.Background(), 10))
	}()

	err := <-output.Errors()

	if err == nil {
		t.Error("Expected an error")
	}
}
