// Package newsletter is ----------------
package newsletter

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
)

// GetReferences returns the content of a url as a string
func GetReferences(ref string) (string, error) {
	resp, err := http.Get(ref)
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
		slog.Warn(fmt.Sprintf("%s returned status code %d", ref, resp.StatusCode))
		return "", nil
	}

	return bodyString, nil
}

// Worker use a worker pool to process jobs and send the restuls through a channel
func Worker(jobs <-chan string, result chan<- string, f func(string) (string, error)) {
	for j := range jobs {
		content, err := f(j)
		if err != nil {
			fmt.Printf("error getting reference %s: %v", j, err)
		}
		result <- content
	}
}
