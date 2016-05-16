package pipeline

import (
	"golang.org/x/net/context"
	"testing"
)

func TestFilter(t *testing.T) {
	var (
		values = make([]context.Context, 10)
	)

	predicate := func(ctx context.Context) bool {
		return FromContext(ctx) > 5
	}

	for i := 0; i < len(values); i++ {
		values[i] = NewContext(context.Background(), i+1)
	}

	outputs, errors := runPipeline(Filter(predicate), values)

	if len(errors) > 0 {
		t.Errorf("Expected %d errors, got %d", 0, len(errors))
	}

	if len(outputs) != 5 {
		t.Errorf("Expected %d values, got %d", 5, len(outputs))
	}
}
