package parser

import (
	"errors"
	"fmt"
	"scraper/log"
	"strings"

	"golang.org/x/net/html"
)

// note that all of this is super brittle

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
	var tableNode *html.Node = nil
	for node := range parse.Descendants() {
		if strings.Contains(node.Data, "table") && strings.Contains(getAttrTraverseSafely(node, "class"), "item-drops") {
			tableNode = node
			break
		}
	}
	if tableNode == nil {
		return nil, errors.New("could not find table tag with item-drops class")
	}

	var tbodyNode *html.Node = nil
	for node := range tableNode.Descendants() {
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
		if getDataTraverseSafely(row, firstChild) == "th" {
			row = row.NextSibling
			continue
		}
		nameLinkTd := row.FirstChild
		quantityTd := traverseSafely(nameLinkTd, nextSibling, nextSibling)
		rarityTd := traverseSafely(quantityTd, nextSibling)
		itemSource := &ItemSource{
			Item: item,
			Source: &Source{
				Name: strings.TrimSuffix(getDataTraverseSafely(nameLinkTd,
					firstChild, // a
					firstChild),
					" "),
				Link: getAttrTraverseSafely(nameLinkTd,
					"href",
					firstChild), // a
				Subclassification: getSubclassification(nameLinkTd),
			},
			Quantity: getDataTraverseSafely(quantityTd, firstChild),
			Rarity: getDataTraverseSafely(rarityTd,
				firstChild, // span
				firstChild),
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
			Name: getDataTraverseSafely(imgNameTd,
				firstChild,  // span (image)
				nextSibling, // " " (space)
				nextSibling, // a (article link)
				firstChild), // item name
			Image: getAttrTraverseSafely(imgNameTd,
				"src",
				firstChild,  // span (image)
				firstChild,  // a (image link)
				firstChild), // img

			Link: getAttrTraverseSafely(imgNameTd,
				"href",
				firstChild,   // span (image)
				nextSibling,  // " " (space)
				nextSibling), // a (article link)
			CompPercent: getDataTraverseSafely(percentTd, firstChild),
		}
		items = append(items, item)
		row = row.NextSibling
	}
	return items, nil
}

func ItemSourceNone(item *Item) *ItemSource {
	return &ItemSource{
		Item: item,
		Source: &Source{
			Name:              "None",
			Link:              "",
			Subclassification: "",
		},
		Quantity: "N/A",
		Rarity:   "N/A",
	}
}

type traversal int

const (
	firstChild traversal = iota
	nextSibling
)

func getAttrTraverseSafely(node *html.Node, key string, path ...traversal) string {
	targetNode := traverseSafely(node, path...)
	if targetNode == nil {
		return ""
	}
	for _, attr := range targetNode.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getDataTraverseSafely(node *html.Node, path ...traversal) string {
	targetNode := traverseSafely(node, path...)
	if targetNode == nil {
		return ""
	}
	return targetNode.Data
}

func getSubclassification(nameLinkTd *html.Node) string {
	nameNode := traverseSafely(nameLinkTd,
		firstChild, // a
		firstChild) // name
	if nameNode == nil {
		return ""
	}
	if nameNode.NextSibling == nil {
		return "" // doing this outside of traverseSafely to avoid the debug line
	}
	return getDataTraverseSafely(nameNode,
		nextSibling, // span
		firstChild)  // subclassification name
}

func traverseSafely(node *html.Node, path ...traversal) *html.Node {
	if node == nil {
		log.Debugln("nil root in traversal")
		return nil
	}
	nextNode := node
	traversed := make([]traversal, 0, len(path))
	for _, next := range path {
		traversed = append(traversed, next)
		switch next {
		case firstChild:
			nextNode = nextNode.FirstChild
		case nextSibling:
			nextNode = nextNode.NextSibling
		default:
			log.Debugf("unknown traversal in path %d\n", next)
			return nil
		}
		if nextNode == nil {
			log.Debugf("nil node in traversal: %s%s\n", node.Data, traversalStr(traversed))
			return nil
		}
	}
	return nextNode
}

func traversalStr(t []traversal) string {
	if len(t) == 0 {
		return ""
	}
	stringBuilder := strings.Builder{}
	for _, tt := range t {
		switch tt {
		case firstChild:
			stringBuilder.WriteString(".firstChild")
		case nextSibling:
			stringBuilder.WriteString(".nextSibling")

		default:
			stringBuilder.WriteString(fmt.Sprintf(".%d", tt))
		}
	}
	return stringBuilder.String()
}
