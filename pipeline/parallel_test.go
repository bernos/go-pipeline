package pipeline

import (
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	concurrency := 5
	delay := time.Millisecond * 500
	expect := time.Millisecond * time.Duration(500*concurrency)

	inner := Map(MapperFunc(func(ctx context.Context) (context.Context, error) {
		time.Sleep(delay)
		x := FromContext(ctx)
		return NewContext(ctx, x+1), nil
	}))

	pl := Parallel(inner, concurrency)
	in, cls := stream.New()
	out := pl(in)

	go func() {
		defer cls()

		for i := 0; i < concurrency; i++ {
			in.Value(NewContext(context.Background(), i))
		}
	}()

	start := time.Now()

	for ctx := range out.Values() {
		_ = FromContext(ctx)
	}

	d := time.Since(start)

	if !(d < expect) {
		t.Errorf("Expected pipeline to finish within %d, but took %d", expect, d)
	}
}
