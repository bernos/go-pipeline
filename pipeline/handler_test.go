package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"testing"
	// "time"
	"fmt"
	"sync"
)

type contextkey int

const (
	valuekey contextkey = iota
)

func NewContext(ctx context.Context, i int) context.Context {
	return context.WithValue(ctx, valuekey, i)
}

func FromContext(ctx context.Context) int {
	if i, ok := ctx.Value(valuekey).(int); ok {
		return i
	}
	return 0
}

func multiply(x int) Handler {
	return HandlerFunc(func(ctx context.Context, out stream.Stream) {
		out.Value(NewContext(ctx, FromContext(ctx)*x))
	})
}

func TestPipe(t *testing.T) {
	want := []int{2, 4, 6, 8, 10}
	input := make([]context.Context, 5)

	for i := 0; i < 5; i++ {
		input[i] = NewContext(context.Background(), i+1)
	}

	values, errors := runPipeline(Pipe(multiply(2)), input)

	if len(errors) != 0 {
		t.Errorf("Expected %d errors, got %d", 0, len(errors))
	}

	if len(values) != len(want) {
		t.Errorf("Expected %d values, got %d", len(want), len(values))
	}

	for i, ctx := range values {
		if FromContext(ctx) != want[i] {
			t.Errorf("Want %d, got %d", want[i], FromContext(ctx))
		}
	}
}

func runPipeline(pl Pipeline, input []context.Context) ([]context.Context, []error) {
	var (
		wg     sync.WaitGroup
		values = make([]context.Context, 0)
		errors = make([]error, 0)
		in     = stream.NewStream()
	)

	out := pl(in)
	wg.Add(2)

	go func() {
		defer wg.Done()
		for value := range out.Values() {
			values = append(values, value)
		}
	}()

	go func() {
		defer wg.Done()
		for err := range out.Errors() {
			errors = append(errors, err)
		}
	}()

	go func() {
		defer in.Close()
		for _, value := range input {
			in.Value(value)
		}

	}()

	wg.Wait()

	return values, errors
}

func TestErrors(t *testing.T) {
	stage := HandlerFunc(func(ctx context.Context, out stream.Stream) {
		out.Error(fmt.Errorf("foo"))
	})

	input := make([]context.Context, 10)

	for i := 0; i < 10; i++ {
		input[i] = NewContext(context.Background(), i)
	}

	values, errors := runPipeline(Pipe(stage), input)

	if len(values) != 0 {
		t.Errorf("Expected %d values, got %d", 0, len(values))
	}

	if len(errors) != 10 {
		t.Errorf("Expected %d errors, got %d", 10, len(errors))
	}
}

func TestCompose(t *testing.T) {
	pl := Compose(Pipe(multiply(2)), Pipe(multiply(3)))

	in := stream.NewStream()
	out := pl(in)

	go func() {
		in.Value(NewContext(context.Background(), 1))
	}()

	go func() {
		for err := range out.Errors() {
			t.Errorf("%s", err.Error())
		}
	}()

	ctx := <-out.Values()

	want := 6
	got := ctx.Value(valuekey).(int)

	if want != got {
		t.Errorf("Want %d, got %d", want, got)
	}
}
