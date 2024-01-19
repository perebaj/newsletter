// Package newsletter is ----------------
package newsletter

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

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

// Worker use a worker pool to process jobs and send the restuls through a channel
func Worker(wg *sync.WaitGroup, urls <-chan string, result chan<- string, f func(string) (string, error)) {
	defer wg.Done()
	for url := range urls {
		content, err := f(url)
		if err != nil {
			slog.Error(fmt.Sprintf("error getting reference: %s", url), "error", err)
		}
		result <- content
	}
}
