package pipeline

import (
	"log"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/bernos/go-pipeline/pipeline/stream"
)

func TestDelay(t *testing.T) {
	p := Delay(time.Millisecond)

	in, cls := stream.New()
	out := p(in)

	go func() {
		defer cls()
		for i := 0; i < 5; i++ {
			in.Value(NewContext(context.Background(), 1))
		}
	}()

	start := time.Now()
	for ctx := range out.Values() {
		log.Printf("%v", ctx)
	}
	duration := time.Since(start)

	if duration < (time.Millisecond * 5) {
		t.Errorf("Expected at least 5 miliseconds, but took %d", duration)
	}
}
