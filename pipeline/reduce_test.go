package pipeline

import (
	"fmt"
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"sync"
	"testing"
)

type IntReducer func(x int, acc int) int

func (r IntReducer) Reduce(ctx context.Context, acc context.Context) (context.Context, error) {
	return NewContext(acc, r(FromContext(ctx), FromContext(acc))), nil
}

func Sum() IntReducer {
	return IntReducer(func(x int, acc int) int {
		return x + acc
	})
}

func TestReduceLeft(t *testing.T) {
	pl := ReduceLeft(Sum())
	input, cls := stream.New()
	output := pl(input)

	go func() {
		defer cls()
		input.Value(NewContext(context.Background(), 1))
		input.Value(NewContext(context.Background(), 2))
		input.Value(NewContext(context.Background(), 3))
		input.Value(NewContext(context.Background(), 4))
		input.Value(NewContext(context.Background(), 5))
	}()

	sum := <-output.Values()
	got := FromContext(sum)

	if got != 15 {
		t.Errorf("Want %d, got %d", 15, got)
	}
}

func TestReduceRight(t *testing.T) {
	pl := ReduceRight(Sum())
	input, cls := stream.New()
	output := pl(input)

	go func() {
		defer cls()
		input.Value(NewContext(context.Background(), 1))
		input.Value(NewContext(context.Background(), 2))
		input.Value(NewContext(context.Background(), 3))
		input.Value(NewContext(context.Background(), 4))
		input.Value(NewContext(context.Background(), 5))
	}()

	sum := <-output.Values()
	got := FromContext(sum)

	if got != 15 {
		t.Errorf("Want %d, got %d", 15, got)
	}
}

func TestReduceLeftError(t *testing.T) {
	var (
		wg       sync.WaitGroup
		sum      int
		errcount int
	)

	pl := ReduceLeft(ReducerFunc(func(ctx context.Context, acc context.Context) (context.Context, error) {
		x := FromContext(ctx)
		y := FromContext(acc)

		if x == 3 {
			return nil, fmt.Errorf("Test error")
		}

		return NewContext(acc, x+y), nil
	}))

	input, cls := stream.New()
	output := pl(input)

	go func() {
		defer cls()
		input.Value(NewContext(context.Background(), 1))
		input.Value(NewContext(context.Background(), 2))
		input.Value(NewContext(context.Background(), 3))
		input.Value(NewContext(context.Background(), 4))
		input.Value(NewContext(context.Background(), 5))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := <-output.Values()
		sum = FromContext(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _ = range output.Errors() {
			errcount++
		}
	}()

	wg.Wait()

	if sum != 12 {
		t.Errorf("Want %d, got %d", 12, sum)
	}

	if errcount != 1 {
		t.Errorf("Wanted 1 error, got %d", errcount)
	}
}

func TestReduceRightError(t *testing.T) {
	var (
		wg       sync.WaitGroup
		sum      int
		errcount int
	)

	pl := ReduceRight(ReducerFunc(func(ctx context.Context, acc context.Context) (context.Context, error) {
		x := FromContext(ctx)
		y := FromContext(acc)

		if x == 3 {
			return nil, fmt.Errorf("Test error")
		}

		return NewContext(acc, x+y), nil
	}))

	input, cls := stream.New()
	output := pl(input)

	go func() {
		defer cls()
		input.Value(NewContext(context.Background(), 1))
		input.Value(NewContext(context.Background(), 2))
		input.Value(NewContext(context.Background(), 3))
		input.Value(NewContext(context.Background(), 4))
		input.Value(NewContext(context.Background(), 5))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := <-output.Values()
		sum = FromContext(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _ = range output.Errors() {
			errcount++
		}
	}()

	wg.Wait()

	if sum != 12 {
		t.Errorf("Want %d, got %d", 12, sum)
	}

	if errcount != 1 {
		t.Errorf("Wanted 1 error, got %d", errcount)
	}
}
