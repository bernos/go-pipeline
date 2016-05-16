package pipeline

import (
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	n := 50

	stage := StageFunc(func(ctx context.Context) (context.Context, error) {
		time.Sleep(time.Second)
		return NewContext(context.Background(), FromContext(ctx)+1), nil
	})

	values := make([]context.Context, n)

	for i := 0; i < n; i++ {
		values[i] = NewContext(context.Background(), i)
	}

	pl := Parallel(Pipe(stage), n)
	start := time.Now().UTC()

	xs, errors := runPipeline(pl, values)

	duration := time.Since(start)

	if len(xs) != n {
		t.Errorf("Expected %d values, got %d", n, len(xs))
	}

	if len(errors) != 0 {
		t.Errorf("Expected %d errors, got %d", 0, len(errors))
	}

	if duration >= time.Second*20 {
		t.Errorf("Want %d, got %d", time.Second, duration)
	}
}
