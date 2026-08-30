package fetcher

import (
	"encoding/json/v2"
	"os"
	"strings"
)

type mockFetcher struct{}

func (m *mockFetcher) FetchPageSection(pageName string, _ string) (string, error) {
	if strings.Contains(pageName, "Collection_log") {
		return m.fetchClogMock()
	}

	return m.fetchWhipMock()
}

func (m *mockFetcher) fetchClogMock() (string, error) {
	dat, err := os.ReadFile("./mock-data/clogtable.json")
	if err != nil {
		return "", err
	}
	response := wikiParseResponse{}
	err = json.Unmarshal(dat, &response)
	if err != nil {
		return "", err
	}
	return response.Parse.Text.Data, nil
}

func (m *mockFetcher) fetchWhipMock() (string, error) {
	dat, err := os.ReadFile("./mock-data/abyssalwhip.json")
	if err != nil {
		return "", err
	}
	response := wikiParseResponse{}
	err = json.Unmarshal(dat, &response)
	if err != nil {
		return "", err
	}
	return response.Parse.Text.Data, nil
}
