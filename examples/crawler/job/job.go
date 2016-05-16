package job

import (
	"golang.org/x/net/context"
)

type key int

const jobKey key = 0

type Job struct {
	URL  string
	Body string
}

func FromContext(ctx context.Context) (Job, bool) {
	job, ok := ctx.Value(jobKey).(Job)
	return job, ok
}

func NewContext(ctx context.Context, job Job) context.Context {
	return context.WithValue(ctx, jobKey, job)
}
