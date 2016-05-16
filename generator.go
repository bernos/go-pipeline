package main

import (
	"golang.org/x/net/context"
)

type Generator interface {
	Generate(ctx context.Context, out chan context.Context, done chan struct{})
}

type GeneratorFunc func(ctx context.Context, out chan context.Context, done chan struct{})

func (fn GeneratorFunc) Generate(ctx context.Context, out chan context.Context, done chan struct{}) {
	fn(ctx, out, done)
}

func Generate(g Generator) Stage {
	return func(in ...<-chan context.Context) []<-chan context.Context {
		var outs []<-chan context.Context

		for i := range in {
			out := make(chan context.Context)
			outs = append(outs, out)

			go func(in <-chan context.Context, out chan context.Context) {
				done := make(chan struct{})
				defer close(done)

				for ctx := range in {
					go func(ctx context.Context, out chan context.Context, done chan struct{}) {
						defer close(out)
						g.Generate(ctx, out, done)
					}(ctx, out, done)
				}
			}(in[i], out)
		}

		return outs
	}
}
