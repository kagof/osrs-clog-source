package parser

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// note that all of this is super brittle; it could easily panic if the wiki changes formats

type Source struct {
	Name              string
	Link              string
	Subclassification string
}

func (s *Source) String() string {
	if s.Subclassification != "" {
		return fmt.Sprintf("%s (%s)", s.Name, s.Subclassification)
	}
	return s.Name
}

type Item struct {
	Name        string
	Image       string
	Link        string
	CompPercent string
}

func (i *Item) String() string {
	return fmt.Sprintf("%s (%s %s %s)",
		i.Name,
		i.Image,
		i.Link,
		i.CompPercent)
}

type ItemSource struct {
	Item     *Item
	Source   *Source
	Quantity string
	Rarity   string
}

func (i *ItemSource) String() string {
	return fmt.Sprintf("%s<->%s (%s, %s)",
		i.Source,
		i.Item.Name,
		i.Quantity,
		i.Rarity)
}

func ParseItemSources(item *Item, data string) ([]*ItemSource, error) {
	parse, err := html.Parse(strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	var tbodyNode *html.Node = nil
	for node := range parse.Descendants() {
		if strings.Contains(node.Data, "tbody") {
			tbodyNode = node
			break
		}
	}
	if tbodyNode == nil {
		return nil, errors.New("could not find tbody tag")
	}
	var itemSources = make([]*ItemSource, 0)
	var row = tbodyNode.FirstChild
	for row != nil {
		if row.FirstChild.Data == "th" {
			row = row.NextSibling
			continue
		}
		nameLinkTd := row.FirstChild
		quantityTd := nameLinkTd.NextSibling.NextSibling
		rarityTd := quantityTd.NextSibling
		itemSource := &ItemSource{
			Item: item,
			Source: &Source{
				Name: strings.TrimSuffix(nameLinkTd.FirstChild. // a
										FirstChild.
										Data,
					" "),
				Link: strings.TrimPrefix(getAttr(nameLinkTd.FirstChild, // a
					"href"), "/w/"),
				Subclassification: getSubclassification(nameLinkTd),
			},
			Quantity: quantityTd.FirstChild.Data,
			Rarity: rarityTd.FirstChild. // span
							FirstChild.
							Data,
		}
		itemSources = append(itemSources, itemSource)
		row = row.NextSibling
	}
	return itemSources, nil
}

func ParseCollectionLogTable(data string) ([]*Item, error) {
	parse, err := html.Parse(strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	var tbodyNode *html.Node = nil
	for node := range parse.Descendants() {
		if strings.Contains(node.Data, "tbody") {
			tbodyNode = node
			break
		}
	}
	if tbodyNode == nil {
		return nil, errors.New("could not find tbody tag")
	}

	var items = make([]*Item, 0, 1800)
	var row = tbodyNode.FirstChild
	for row != nil {
		if row.FirstChild.Data == "th" {
			row = row.NextSibling
			continue
		}
		imgNameTd := row.FirstChild
		percentTd := imgNameTd.NextSibling.NextSibling

		item := &Item{
			Name: imgNameTd.
				FirstChild.  // span (image)
				NextSibling. // " " (space)
				NextSibling. // a (article link)
				FirstChild.  // item name
				Data,
			Image: getAttr(imgNameTd.
				FirstChild. // span (image)
				FirstChild. // a (image link)
				FirstChild, // img
				"src"),
			Link: strings.TrimPrefix(getAttr(imgNameTd.
				FirstChild.  // span (image)
				NextSibling. // " " (space)
				NextSibling, // a (article link)
				"href"), "/w/"),
			CompPercent: percentTd.FirstChild.Data,
		}
		items = append(items, item)
		row = row.NextSibling
	}
	return items, nil
}

func getAttr(n *html.Node, key string) string {
	if n == nil {
		return ""
	}
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getSubclassification(nameLinkTd *html.Node) string {
	nameNode := nameLinkTd.FirstChild. // a
						FirstChild // name
	if nameNode.NextSibling == nil {
		return ""
	}
	return nameNode.NextSibling.FirstChild.Data
}
