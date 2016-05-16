package pipeline

import (
	"golang.org/x/net/context"
)

// Stage represents an actual stage in a pipeline. It takes values from
// a context, operates on them, and produces a new context. Stages are
// wrapped in Pipeline functions which take care of pushing and pulling
// contexts to and from Stages
type Stage interface {
	Handle(context.Context) (context.Context, error)
}

// StageFunc makes a regular func implement the Stage interface
type StageFunc func(ctx context.Context) (context.Context, error)

// Handle satisfies the Stage interface for StageFunc
func (fn StageFunc) Handle(ctx context.Context) (context.Context, error) {
	return fn(ctx)
}

// Higher order Pipeline funcs

// func Split(pipelines ...Pipeline) Pipeline {
// 	return func(in <-chan context.Context) <-chan context.Context {
// 		out := make(chan context.Context)
// 		pipelineInputs := make([]chan context.Context, len(pipelines))

// 		for i, pipeline := range pipelines {
// 			pipelineInputs[i] = make(chan context.Context)

// 			// Collect outputs from pipelines and re-send on out
// 			go func(pipelineOutput <-chan context.Context) {
// 				for ctx := range pipelineOutput {
// 					out <- ctx
// 				}
// 			}(pipeline(pipelineInputs[i]))
// 		}

// 		go func() {
// 			defer close(out)

// 			for i := range pipelineInputs {
// 				defer close(pipelineInputs[i])
// 			}

// 			for ctx := range in {
// 				for i := range pipelineInputs {
// 					// TODO could wrap this in a go routine, so we dont
// 					// block if one pipeline is slow
// 					pipelineInputs[i] <- ctx
// 				}
// 			}
// 		}()

// 		return out
// 	}
// }
