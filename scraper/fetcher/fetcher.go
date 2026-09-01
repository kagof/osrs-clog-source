package fetcher

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"scraper/cache"
	"scraper/log"
	"strings"
	"time"
)

const (
	BaseUrl   = "https://oldschool.runescape.wiki"
	userAgent = "kagof-clog-scraper/0.0 (osrs-clog-source.kagof.com scraper; https://github.com/kagof/osrs-clog-source)"
)

func New(useMock bool, delay time.Duration) Fetcher {
	if useMock {
		return &mockFetcher{}
	}
	return getFetcher(delay, 15)
}

type Fetcher interface {
	FetchPageSection(pageName string, sectionName string) (string, error)
}

type fetcher struct {
	client *http.Client
	cache  *cache.Cache[string, string]
}

func getFetcher(delay time.Duration, cacheCap int) *fetcher {
	return &fetcher{
		client: &http.Client{
			Transport: newRoundTripper(userAgent, delay),
		},
		cache: cache.NewCache[string, string](cacheCap),
	}
}

type roundTripper struct {
	userAgent string
	delay     time.Duration
	lastTime  time.Time
}

func newRoundTripper(userAgent string, delay time.Duration) *roundTripper {
	return &roundTripper{
		userAgent: userAgent,
		delay:     delay,
		lastTime:  time.Now().Add(-delay),
	}
}

func (u *roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", u.userAgent)
	r.Header.Set("Accept", "application/json")

	// ensures at least u.delay time has passed since the last request was sent.
	now := time.Now()
	waitUntil := u.lastTime.Add(u.delay)
	if waitUntil.After(now) {
		time.Sleep(waitUntil.Sub(now))
	}

	u.lastTime = time.Now()

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}
	return resp, err
}

type wikiParseResponse struct {
	Parse parse `json:"parse"`
}

type parse struct {
	Title    string    `json:"title"`
	Text     text      `json:"text"`
	Sections []section `json:"sections"`
}
type text struct {
	Data string `json:"*"`
}

type section struct {
	Title string `json:"line"`
	Index string `json:"index"`
}

func (f *fetcher) FetchPageSection(pageName string, sectionName string) (string, error) {
	pageName = strings.SplitN(strings.TrimPrefix(pageName, "/w/"), "#", 2)[0]
	key := fmt.Sprintf("%s#%s", pageName, sectionName)
	return f.cache.Memoized(f.fetchPageSectionInternal)(key)
}

func (f *fetcher) fetchPageSectionInternal(key string) (string, error) {
	spl := strings.SplitN(key, "#", 2)
	pageName := spl[0]
	sectionName := spl[1]
	sections, err := f.fetchSections(pageName)
	if err != nil {
		return "", err
	}
	secIdx := ""
	for _, sec := range sections {
		if strings.EqualFold(sec.Title, sectionName) {
			secIdx = sec.Index
			break
		}
	}
	if secIdx == "" {
		return "", fmt.Errorf("%s section not found", sectionName)
	}
	resp, err := f.client.Get(fmt.Sprintf("%s/api.php?action=parse&page=%s&section=%s&format=json&redirects=true",
		BaseUrl,
		pageName,
		secIdx))
	if err != nil {
		return "", fmt.Errorf("error fetching %s section %s (%s), %w", pageName, sectionName, secIdx, err)
	}
	defer func(Body io.ReadCloser) {
		err2 := Body.Close()
		if err2 != nil {
			log.Printf("could not close body: %s", err2)
		}
	}(resp.Body)
	response := wikiParseResponse{}
	err = json.UnmarshalRead(resp.Body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing %s section %s (%s), %w", pageName, sectionName, secIdx, err)
	}
	return response.Parse.Text.Data, nil
}

func (f *fetcher) fetchSections(pageName string) ([]section, error) {
	resp, err := f.client.Get(fmt.Sprintf("%s/api.php?action=parse&page=%s&props=sections&format=json&redirects=true",
		BaseUrl,
		pageName))
	if err != nil {
		return nil, fmt.Errorf("error fetching sections for %s, %w", pageName, err)
	}
	defer func(Body io.ReadCloser) {
		err2 := Body.Close()
		if err2 != nil {
			log.Printf("could not close body: %s", err2)
		}
	}(resp.Body)
	response := wikiParseResponse{}
	err = json.UnmarshalRead(resp.Body, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing sections for %s, %w", pageName, err)
	}
	return response.Parse.Sections, nil
}
