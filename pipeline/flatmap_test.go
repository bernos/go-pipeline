package pipeline

import (
	"fmt"
	"github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"testing"
	"time"
)

type IntFlatMapper func(x int) []int

func (m IntFlatMapper) FlatMap(ctx context.Context) ([]context.Context, error) {
	ints := m(FromContext(ctx))
	ret := make([]context.Context, len(ints))

	for i := range ints {
		ret[i] = NewContext(context.Background(), ints[i])
	}

	return ret, nil
}

func TestFlatMap(t *testing.T) {
	pl := FlatMap(IntFlatMapper(func(x int) []int {
		return []int{x, x + 1, x + 2, x + 3}
	}))

	in, cls := stream.New()
	output := pl(in)

	go func() {
		in.Value(NewContext(context.Background(), 1))
		in.Value(NewContext(context.Background(), 2))
		cls()
	}()

	want := []int{1, 2, 3, 4, 2, 3, 4, 5}
	got := make([]int, 0)

	for ctx := range output.Values() {
		got = append(got, FromContext(ctx))
	}

	if len(got) != len(want) {
		t.Errorf("Want %d items, got %d", len(want), len(got))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Want %d, got %d", want[i], got[i])
		}
	}
}

func TestFlatMapError(t *testing.T) {
	mapper := FlatMapperFunc(func(ctx context.Context) ([]context.Context, error) {
		return nil, fmt.Errorf("An error")
	})

	pl := FlatMap(mapper)

	in, cls := stream.New()
	output := pl(in)

	go func() {
		in.Value(NewContext(context.Background(), 1))
		cls()
	}()

	err := <-output.Errors()

	if err == nil {
		t.Errorf("Expected and error but got none")
	}
}

func TestPFlatMap(t *testing.T) {
	pl := PFlatMap(IntFlatMapper(func(x int) []int {
		time.Sleep(time.Millisecond * 500)
		return []int{x, x + 1, x + 2, x + 3}
	}), 2)

	in, cls := stream.New()
	output := pl(in)

	go func() {
		in.Value(NewContext(context.Background(), 1))
		in.Value(NewContext(context.Background(), 2))
		cls()
	}()

	want := []int{1, 2, 3, 4, 2, 3, 4, 5}
	got := make([]int, 0)

	start := time.Now()
	for ctx := range output.Values() {
		got = append(got, FromContext(ctx))
	}

	d := time.Since(start)

	if len(got) != len(want) {
		t.Errorf("Want %d items, got %d", len(want), len(got))
	}

	if d > time.Millisecond*550 {
		t.Errorf("Expected no longer that %dms, but took %d", 550, d)
	}
}
