// Package newsletter is ----------------
package newsletter

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/perebaj/newsletter/mongodb"
)

// PageContent is the struct that gather important information of a website
type PageContent struct {
	Content string
	URL     string
}

// Storage is the interface that wraps the basic methods to save and get data from the database
type Storage interface {
	SaveSite(ctx context.Context, site []mongodb.Site) error
	DistinctEngineerURLs(ctx context.Context) ([]interface{}, error)
}

// Crawler contains the necessary information to run the crawler
type Crawler struct {
	URLch    chan string
	resultCh chan PageContent
	signalCh chan os.Signal
	MaxJobs  int
	wg       *sync.WaitGroup
	// scheduler is the pace time between each fetch
	scheduler time.Duration
}

// NewCrawler initializes a new Crawler
func NewCrawler(maxJobs int, s time.Duration, signalCh chan os.Signal) *Crawler {
	return &Crawler{
		URLch:     make(chan string),
		resultCh:  make(chan PageContent),
		signalCh:  signalCh,
		wg:        &sync.WaitGroup{},
		MaxJobs:   maxJobs,
		scheduler: s,
	}
}

// Run starts the crawler, where s represents the storage and f the function to fetch the content of a website
func (c *Crawler) Run(ctx context.Context, s Storage, f func(string) (string, error)) {
	c.wg.Add(c.MaxJobs)
	for i := 0; i < c.MaxJobs; i++ {
		go c.Worker(f)
	}

	go func() {
		defer close(c.URLch)
		for range time.Tick(c.scheduler) {
			slog.Debug("fetching engineers")
			gotURLs, err := s.DistinctEngineerURLs(ctx)
			if err != nil {
				slog.Error("error getting engineers", "error", err)
				c.signalCh <- syscall.SIGTERM
			}

			slog.Debug("fetched engineers", "engineers", len(gotURLs))
			for _, url := range gotURLs {
				c.URLch <- url.(string)
			}
		}
	}()

	go func() {
		c.wg.Wait()
		defer close(c.resultCh)
	}()

	go func() {
		for v := range c.resultCh {
			slog.Debug("saving fetched sites response")
			err := s.SaveSite(ctx, []mongodb.Site{
				{
					URL:            v.URL,
					Content:        v.Content,
					ScrapeDatetime: time.Now().UTC(),
				},
			})
			if err != nil {
				slog.Error("error saving site result", "error", err)
				c.signalCh <- syscall.SIGTERM
			}
		}
	}()
}

// Worker use a worker pool to process jobs and send the restuls through a channel
func (c *Crawler) Worker(f func(string) (string, error)) {
	defer c.wg.Done()
	for url := range c.URLch {
		content, err := f(url)
		if err != nil {
			slog.Error(fmt.Sprintf("error getting reference: %s", url), "error", err)
		}
		c.resultCh <- PageContent{Content: content, URL: url}
	}
}

// Fetch returns the content of a url as a string
func Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var bodyString string
	if resp.StatusCode == 200 {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString = buf.String()
	} else {
		slog.Warn(fmt.Sprintf("%s returned status code %d", url, resp.StatusCode))
		return "", nil
	}

	return bodyString, nil
}
