package newsletter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestWorker(_ *testing.T) {
	urls := make(chan string)
	results := make(chan string)

	f := func(s string) (string, error) {
		time.Sleep(100 * time.Millisecond)
		return fmt.Sprintf("job %s done", s), nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go Worker(&wg, urls, results, f)
	go Worker(&wg, urls, results, f)

	go func() {
		urls <- "job1"
		urls <- "job2"
		urls <- "job3"
		urls <- "job4"
		urls <- "job5"
		urls <- "job6"
		defer close(urls)
	}()

	go func() {
		wg.Wait()
		defer close(results)
	}()

	for i := 0; i < 6; i++ {
		<-results
	}
}

func TestGetReferences(t *testing.T) {
	wantBody := "Hello, World!"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(wantBody))
	}))

	defer server.Close()

	got, err := Fetch(server.URL)
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

	got, err := Fetch(server.URL)
	if err != nil {
		t.Errorf("error getting reference: %v", err)
	}

	if got != "" {
		t.Errorf("expected empty body, got %s", got)
	}
}
