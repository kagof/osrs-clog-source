package fetcher

import (
	"os"
	"scraper/log"
	"testing"
	"time"
)

func TestFetcherMemoizationForAncientPages(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration test")
	}
	log.Writer = t.Output()
	log.EnableDebug()
	f := getFetcher(time.Second, 15)

	_, err := f.FetchPageSection("Ancient_page#1", "Item sources")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.FetchPageSection("Ancient_page#2", "Item sources")
	if err != nil {
		t.Fatal(err)
	}
	if f.cache.Len() != 1 {
		t.Fatalf("cache size expected 1, actual %d", f.cache.Len())
	}
}
