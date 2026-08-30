package results

import (
	"fmt"
	"scraper/fetcher"
	"scraper/parser"
	"time"
)

type CollectionLogItemSources struct {
	LastUpdated time.Time                           `json:"lastUpdated"`
	Sources     map[string]*CollectionLogItemSource `json:"sources"`
}

type CollectionLogItemSource struct {
	Name               string                                         `json:"name"`
	Subclassifications map[string]*CollectionLogItemSubclassification `json:"subclassifications"`
}

type CollectionLogItemSubclassification struct {
	Name  string               `json:"name"`
	Link  string               `json:"link"`
	Items []*CollectionLogItem `json:"items"`
}

type CollectionLogItem struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Link        string `json:"link"`
	CompPercent string `json:"compPercent"`
	Quantity    string `json:"quantity"`
	Rarity      string `json:"rarity"`
}

func New() *CollectionLogItemSources {
	return &CollectionLogItemSources{
		LastUpdated: time.Now().Round(time.Second).UTC(),
		Sources:     make(map[string]*CollectionLogItemSource),
	}
}

func (s *CollectionLogItemSources) Add(is *parser.ItemSource) {
	source, ok := s.Sources[is.Source.Name]
	if !ok {
		source = &CollectionLogItemSource{
			Name:               is.Source.Name,
			Subclassifications: make(map[string]*CollectionLogItemSubclassification),
		}
		s.Sources[is.Source.Name] = source
	}
	subclassification, ok := source.Subclassifications[is.Source.Subclassification]
	if !ok {
		subclassification = &CollectionLogItemSubclassification{
			Name:  is.Source.Subclassification,
			Link:  buildPageLink(is.Source.Link),
			Items: make([]*CollectionLogItem, 0),
		}
		source.Subclassifications[is.Source.Subclassification] = subclassification
	}
	subclassification.Items = append(subclassification.Items, &CollectionLogItem{
		Name:        is.Item.Name,
		Image:       buildImageLink(is.Item.Image),
		Link:        buildPageLink(is.Item.Link),
		CompPercent: is.Item.CompPercent,
		Quantity:    is.Quantity,
		Rarity:      is.Rarity,
	})
}

func buildPageLink(link string) string {
	return fmt.Sprintf("%s/w/%s", fetcher.BaseUrl, link)
}

func buildImageLink(link string) string {
	return fmt.Sprintf("%s%s", fetcher.BaseUrl, link)
}
