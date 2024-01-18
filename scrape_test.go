package newsletter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	const numJobs = 3
	jobs := make(chan string, numJobs)
	result := make(chan string, numJobs)

	f := func(s string) (string, error) {
		time.Sleep(100 * time.Millisecond)
		return fmt.Sprintf("job %s done", s), nil
	}

	go Worker(jobs, result, f)

	jobs <- "job1"
	jobs <- "job2"
	jobs <- "job3"

	close(jobs)

	for i := 0; i < numJobs; i++ {
		r := <-result
		if r != fmt.Sprintf("job job%d done", i+1) {
			t.Errorf("expected job%d done, got %s", i+1, r)
		}
	}
}

func TestGetReferences(t *testing.T) {
	wantBody := "Hello, World!"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(wantBody))
	}))

	defer server.Close()

	got, err := GetReferences(server.URL)
	if err != nil {
		t.Errorf("error getting reference: %v", err)
	}

	if got != wantBody {
		t.Errorf("expected %s, got %s", wantBody, got)
	}
}

func TestGetReferences_Status500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	got, err := GetReferences(server.URL)
	if err != nil {
		t.Errorf("error getting reference: %v", err)
	}

	if got != "" {
		t.Errorf("expected empty body, got %s", got)
	}
}
