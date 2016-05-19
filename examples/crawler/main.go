package main

import (
	"github.com/bernos/go-pipeline/examples/crawler/job"
	"github.com/bernos/go-pipeline/pipeline"
	// "github.com/bernos/go-pipeline/pipeline/stream"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

var (
	urlRegexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
)

func main() {

	crawler := pipeline.
		PMap(fetchURLMap(&http.Client{}), 20).
		Map(saveFileMap()).
		FlatMap(findURLSMap()).
		Filter(dedupe())

	ctx, _ := context.WithTimeout(job.NewContext(context.Background(), job.Job{URL: "http://www.wikipedia.com"}), time.Second*15)

	out := crawler.Loop(ctx)

	for ctx := range out.Values() {
		j, _ := job.FromContext(ctx)
		log.Printf("Finished one - %s\n", j.URL)
	}

	log.Println("Done!")
}

func dedupe() pipeline.Predicate {
	history := make(map[string]bool)

	return func(ctx context.Context) bool {
		seen := false

		if j, ok := job.FromContext(ctx); ok {
			seen = history[j.URL]
			// log.Printf("Seen %s - %t", j.URL, seen)
			history[j.URL] = true
		}

		return !seen
	}
}

func fetchURLMap(client *http.Client) pipeline.Mapper {
	return job.Mapper(func(j job.Job) (job.Job, error) {
		log.Printf("fetching %s\n", j.URL)

		resp, err := client.Get(j.URL)

		if err != nil {
			return j, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return j, err
		}

		j.Body = string(body)

		return j, nil
	})
}

func saveFileMap() pipeline.Mapper {
	return job.Mapper(func(j job.Job) (job.Job, error) {
		log.Printf("Saving %s\n", j.URL)
		return j, nil
	})
}

func counter() pipeline.Mapper {
	count := 0
	return job.Mapper(func(j job.Job) (job.Job, error) {
		count++
		log.Printf("Counted %d\n", count)
		return j, nil
	})
}

func findURLSMap() pipeline.FlatMapper {
	return job.FlatMapper(func(j job.Job) ([]job.Job, error) {
		result := urlRegexp.FindAllString(j.Body, -1)
		jobs := make([]job.Job, len(result))
		log.Printf("Found %d urls", len(result))
		// log.Printf("%v\n", result)

		if result != nil {
			for i, url := range result {
				jobs[i] = job.Job{URL: url}
			}
		}

		return jobs, nil
	})
}

func fetchURL(client *http.Client) pipeline.Handler {
	return job.Handler(func(j job.Job, out func(job.Job) error) error {
		log.Printf("fetching %s\n", j.URL)

		resp, err := client.Get(j.URL)

		if err != nil {
			return err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		j.Body = string(body)

		return out(j)
	})
}

func saveFile() pipeline.Handler {
	return job.Handler(func(j job.Job, out func(job.Job) error) error {
		log.Printf("Saving %s\n", j.URL)
		return out(j)
	})
}

func findURLS() pipeline.Handler {
	return job.Handler(func(j job.Job, out func(job.Job) error) error {
		result := urlRegexp.FindAllString(j.Body, -1)
		log.Printf("Found %d urls", len(result))
		if result != nil {
			for _, url := range result {
				out(job.Job{URL: url})
			}
		}
		return nil
	})
}
