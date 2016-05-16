package pipeline

import (
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

func multiply(x int) Stage {
	return StageFunc(func(ctx context.Context) (context.Context, error) {
		return NewContext(ctx, FromContext(ctx)*x), nil
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
		in     = make(chan context.Context)
	)

	out, errs := pl(in)
	wg.Add(2)

	go func() {
		defer close(in)
		for _, value := range input {
			in <- value
		}
	}()

	go func() {
		defer wg.Done()
		for value := range out {
			values = append(values, value)
		}
	}()

	go func() {
		defer wg.Done()
		for err := range errs {
			errors = append(errors, err)
		}
	}()

	wg.Wait()

	return values, errors
}

func TestErrors(t *testing.T) {
	stage := StageFunc(func(ctx context.Context) (context.Context, error) {
		return nil, fmt.Errorf("foo")
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

// func TestParallelPipe(t *testing.T) {
// 	n := 10
// 	in := make(chan context.Context)

// 	go func() {
// 		defer close(in)
// 		for i := 0; i < n; i++ {
// 			in <- NewContext(context.Background(), 1)
// 		}
// 	}()

// 	// Make a stage that sleeps for a second and then returns some value.
// 	// Running n of these stages in parallel should not take longer than
// 	// running one.
// 	stage := StageFunc(func(ctx context.Context) context.Context {
// 		time.Sleep(time.Second * 1)
// 		return NewContext(ctx, FromContext(ctx)+1)
// 	})

// 	pl := ParallelPipe(stage, n)
// 	out := pl(in)
// 	start := time.Now().UTC()

// 	<-out

// 	duration := time.Since(start)

// 	if duration-time.Second > time.Millisecond*50 {
// 		t.Errorf("Want %d, got %d", time.Second, duration)
// 	}
// }

func TestCompose(t *testing.T) {
	pl := Compose(Pipe(multiply(2)), Pipe(multiply(3)))

	in := make(chan context.Context)
	out, errors := pl(in)

	go func() {
		in <- NewContext(context.Background(), 1)
	}()

	go func() {
		for err := range errors {
			t.Errorf("%s", err.Error())
		}
	}()

	ctx := <-out

	want := 6
	got := ctx.Value(valuekey).(int)

	if want != got {
		t.Errorf("Want %d, got %d", want, got)
	}
}
