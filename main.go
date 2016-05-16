package main

import (
	"github.com/bernos/go-pipeline/pipeline"
	"golang.org/x/net/context"
)

func main() {
	pl := pipeline.ParallelPipe(fetchURL, 10).Compose()
}

// Fetches content from a url and returns it, wrapped in a context
func fetchURL(ctx context.Context) context.Context {
	return context.Background()
}

func saveFile(ctx context.Context) context.Context {
	return context.Background()
}

func findURLS(ctx context.Context) context.Context {
	return context.Background()
}
