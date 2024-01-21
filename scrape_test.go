package newsletter

import (
	"context"
	"crypto/md5"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/perebaj/newsletter/mongodb"
)

const fakeURL = "http://fakeurl.test"

func TestPageComparation(t *testing.T) {
	recentScrapedPage := Page{
		Content:        "Hello, World!",
		URL:            fakeURL,
		ScrapeDateTime: time.Now().UTC(),
	}

	lastScrapedPage := []mongodb.Page{
		{
			Content:        "Hello, World!",
			URL:            fakeURL,
			ScrapeDatetime: time.Now().UTC().Add(-time.Duration(1) * time.Hour),
			HashMD5:        md5.Sum([]byte("Hello, World!")),
		},
	}

	newPage := pageComparation(lastScrapedPage, recentScrapedPage)

	if newPage[0].IsMostRecent {
		t.Errorf("expected false, got %v", newPage[0].IsMostRecent)
	}

	lastScrapedPage[0].Content = "Hello, World! 2"
	lastScrapedPage[0].HashMD5 = md5.Sum([]byte("Hello, World! 2"))

	newPage = pageComparation(lastScrapedPage, recentScrapedPage)

	if !newPage[0].IsMostRecent {
		t.Errorf("expected true, got %v", newPage[0].IsMostRecent)
	}

	lastScrapedPage = []mongodb.Page{}

	newPage = pageComparation(lastScrapedPage, recentScrapedPage)

	if !newPage[0].IsMostRecent {
		t.Errorf("expected true, got %v", newPage[0].IsMostRecent)
	}
}

// Even not verifying the result, this test is useful to check if the crawler is running properly, since it is
// using Mocks for the Storage and the Fetch function.
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

type StorageMockImpl struct {
}

func NewStorageMock() StorageMockImpl {
	return StorageMockImpl{}
}

func (s StorageMockImpl) SavePage(_ context.Context, _ []mongodb.Page) error {
	return nil
}

func (s StorageMockImpl) DistinctEngineerURLs(_ context.Context) ([]interface{}, error) {
	return []interface{}{fakeURL}, nil
}

func (s StorageMockImpl) Page(_ context.Context, _ string) ([]mongodb.Page, error) {
	return []mongodb.Page{}, nil
}
