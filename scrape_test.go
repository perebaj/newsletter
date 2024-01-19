package newsletter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/perebaj/newsletter/mongodb"
)

const fakeURL = "http://fakeurl.test"

func TestCrawlerRun(t *testing.T) {
	timeoutCh := time.After(time.Duration(150) * time.Millisecond)
	ctx := context.Background()
	s := NewStorageMock()

	f := func(string) (string, error) {
		return "Hello, World!", nil
	}

	signalCh := make(chan os.Signal, 1)

	c := NewCrawler(1, time.Duration(1000)*time.Millisecond, signalCh)
	go func() {
		c.Run(ctx, s, f)
	}()

	select {
	case <-signalCh:
		t.Error("unexpected signal error")
	case <-timeoutCh:
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

type StorageMock interface {
	SaveSite(ctx context.Context, site []mongodb.Site) error
	DistinctEngineerURLs(ctx context.Context) ([]interface{}, error)
}

type StorageMockImpl struct {
}

func NewStorageMock() StorageMock {
	return StorageMockImpl{}
}

func (s StorageMockImpl) SaveSite(ctx context.Context, site []mongodb.Site) error {
	return nil
}

func (s StorageMockImpl) DistinctEngineerURLs(ctx context.Context) ([]interface{}, error) {
	return []interface{}{fakeURL}, nil
}
