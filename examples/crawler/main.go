package main

import (
	"github.com/bernos/go-pipeline/examples/crawler/job"
	"github.com/bernos/go-pipeline/pipeline"
	"github.com/bernos/go-pipeline/pipeline/stream"
	// "github.com/bernos/go-pipeline/pipeline/stream"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/context"
)

var (
	urlRegexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
)

func main() {

	// A webcrawler pipeline that will recursively crawl a website, downloading content
	// in parallel, and removing duplicate urls
	crawler := pipeline.
		PMap(fetchURL(&http.Client{}), 10).
		Map(saveFile()).
		FlatMap(findURLS()).
		Filter(dedupe()).
		Loop()

	// Point the crawler at wikipedia, and configure a timeout using the context
	ctx, cancel := context.WithTimeout(job.NewContext(context.Background(), job.Job{URL: "http://www.wikipedia.com"}), time.Second*15)

	in, cls := stream.New()
	out := crawler(in)
	in.Value(ctx)

	// Start a go routine to monitor for errors on the pipeline error channel.
	// For now we will just stop the pipeline, using the cancel func for our
	// context
	go func() {
		for err := range out.Errors() {
			log.Printf("Error: %s\n", err.Error())
			cancel()
		}
	}()

	// Close the input stream when the context times out
	go func() {
		<-ctx.Done()
		cls()
	}()

	// Print out the pipeline output
	for ctx := range out.Values() {
		j, _ := job.FromContext(ctx)
		log.Printf("Found a link to %s\n", j.URL)
	}

	log.Println("Done!")
}

// dedupe returns a pipeline Predicate that returns false if we have seen the
// URL to be crawled before
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

// fetchURL returns a pipeline Mapper that fetches the content for a URL and
// adds it to the job in the context
func fetchURL(client *http.Client) pipeline.Mapper {
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

// saveFile returns a pipeline Mapper that saves the content downloaded from
// a url to disk.
// TODO: implement the saving bit
func saveFile() pipeline.Mapper {
	return job.Mapper(func(j job.Job) (job.Job, error) {
		log.Printf("Saving %s\n", j.URL)
		return j, nil
	})
}

// findURLS returns a pipeline FlatMapper that scans the content downloaded from
// a urls for any URLs that appear in it
func findURLS() pipeline.FlatMapper {
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
