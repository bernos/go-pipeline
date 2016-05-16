package main

import (
	"fmt"
	"github.com/bernos/go-pipeline/examples/crawler/job"
	"github.com/bernos/go-pipeline/pipeline"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func main() {
	input := make(chan context.Context)

	defer close(input)

	crawler := pipeline.
		ParallelPipe(fetchURL(&http.Client{}), 10).
		Pipe(saveFile()).
		Pipe(findURLS()).
		Filter(dedupe())

	out, _ := crawler(input)
	ctx, _ := context.WithTimeout(job.NewContext(context.Background(), job.Job{URL: "http://www.wikipedia.com"}), time.Second*50)
	done := ctx.Done()

	input <- ctx

	for {
		select {
		case <-done:
			log.Printf("Finished!")
			return
		case ctx := <-out:
			go func(ctx context.Context) {
				input <- ctx
			}(ctx)
		}
	}
}

func dedupe() pipeline.Predicate {
	history := make(map[string]bool)

	return func(ctx context.Context) bool {
		seen := false

		if j, ok := job.FromContext(ctx); ok {
			seen = history[j.URL]
			history[j.URL] = true
		}

		return !seen
	}
}

func jobStage(fn func(job.Job, chan context.Context, chan error)) pipeline.Stage {
	return pipeline.StageFunc(func(ctx context.Context, out chan context.Context, errors chan error) {
		if j, ok := job.FromContext(ctx); ok {
			fn(j, out, errors)
		} else {
			errors <- fmt.Errorf("Unable to find job in context")
		}
	})
}

func fetchURL(client *http.Client) pipeline.Stage {
	return pipeline.StageFunc(func(ctx context.Context, out chan context.Context, errors chan error) {
		if j, ok := job.FromContext(ctx); ok {
			log.Printf("fetching %s\n", j.URL)

			resp, err := client.Get(j.URL)

			if err != nil {
				errors <- err
				return
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				errors <- err
				return
			}

			j.Body = string(body)
			out <- job.NewContext(ctx, j)
		} else {
			errors <- fmt.Errorf("Unable to find job in context")
		}
	})
}

func saveFile() pipeline.Stage {
	return pipeline.StageFunc(func(ctx context.Context, out chan context.Context, errors chan error) {
		if j, ok := job.FromContext(ctx); ok {
			log.Printf("Saving %s\n", j.URL)

			out <- ctx
		} else {
			errors <- fmt.Errorf("Unable to find job in context")
		}
	})
}

func findURLS() pipeline.Stage {
	re := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	return pipeline.StageFunc(func(ctx context.Context, out chan context.Context, errors chan error) {
		if j, ok := job.FromContext(ctx); ok {
			result := re.FindAllString(j.Body, -1)
			log.Printf("Found %d urls", len(result))
			if result != nil {
				for _, url := range result {
					out <- job.NewContext(context.Background(), job.Job{URL: url})
				}
			}
		} else {
			errors <- fmt.Errorf("Unable to find job in context")
		}
	})
}
