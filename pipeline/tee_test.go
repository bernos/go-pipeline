package pipeline

import (
	"fmt"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func SkipTestTee(t *testing.T) {
	var (
		n         = 10
		collected = make([]int, 0)
		inputs    = make([]context.Context, n)
	)

	sink := func(ctx context.Context) error {
		time.Sleep(time.Millisecond)
		fmt.Println("SINK")
		collected = append(collected, FromContext(ctx))
		return nil
	}

	for i := 0; i < n; i++ {
		inputs[i] = NewContext(context.Background(), i+1)
	}

	values, errors := runPipeline(Tee(Sink(sink)), inputs)

	if len(errors) != 0 {
		t.Errorf("Expected %d errors, got %d", 0, len(errors))
	}

	if len(collected) != n {
		t.Errorf("Expected collected length to be %d, but got %d", n, len(collected))
	}

	if len(values) != n {
		t.Errorf("Expected values length to be %d, but got %d", n, len(values))
	}
}
