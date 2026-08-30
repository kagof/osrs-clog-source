package main

import (
	"encoding/json"
	"fmt"
	"os"
	"scraper/fetcher"
	"scraper/log"
	"scraper/parser"
	"scraper/results"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		useMock         = false
		pretty          = false
		tee             = false
		outFile         = ""
		truncateItemsTo = -1
	)
	if slices.Contains(os.Args, "--mock") {
		useMock = true
		log.Println("Using mock fetcher")
	}
	if slices.Contains(os.Args, "--pretty") {
		pretty = true
		log.Println("Pretty-printing")
	}
	outputIndex := slices.IndexFunc(os.Args, func(s string) bool {
		return strings.HasPrefix(s, "-o=")
	})
	if outputIndex >= 0 {
		o := os.Args[outputIndex]
		outFile = strings.SplitN(o, "=", 2)[1]
		log.Printf("Using output file %s\n", outFile)
		if slices.Contains(os.Args, "--tee") {
			tee = true
			log.Println("Also writing output to stdout")
		}
	}
	truncateIndex := slices.IndexFunc(os.Args, func(s string) bool {
		return strings.HasPrefix(s, "-t=")
	})
	if truncateIndex >= 0 {
		t := os.Args[truncateIndex]
		tVal := strings.SplitN(t, "=", 2)[1]
		tInt, err := strconv.Atoi(tVal)
		if err != nil {
			panic(err)
		}
		truncateItemsTo = tInt
		log.Printf("Truncating items to %d\n", truncateItemsTo)
	}
	Execute(useMock, pretty, tee, outFile, truncateItemsTo)
}

func Execute(useMock, pretty, tee bool, outFile string, truncateItemsTo int) {
	startTime := time.Now()
	f := fetcher.New(useMock, time.Second)

	clogTable, err := f.FetchPageSection("Collection_log%2FTable", "Table")
	if err != nil {
		panic(err)
	}
	items, err := parser.ParseCollectionLogTable(clogTable)
	if err != nil {
		panic(err)
	}
	if truncateItemsTo >= 0 {
		items = items[:truncateItemsTo]
	}

	resultSet := results.New()
	for _, item := range items {
		sec, err2 := f.FetchPageSection(item.Link, "Item sources")
		if err2 != nil {
			panic(err2)
		}
		sources, err2 := parser.ParseItemSources(item, sec)
		if err2 != nil {
			panic(err2)
		}
		for _, source := range sources {
			resultSet.Add(source)
		}
	}
	var bytes []byte
	if pretty {
		bytes, err = json.MarshalIndent(resultSet, "", "  ")
		if err != nil {
			panic(err)
		}
	} else {
		bytes, err = json.Marshal(resultSet)
		if err != nil {
			panic(err)
		}
	}
	log.Printf("Scraped %d items into %d sources in %v\n", len(items), len(resultSet.Sources), time.Since(startTime))

	if outFile != "" {
		err2 := os.WriteFile(outFile, bytes, 0644)
		if err2 != nil {
			panic(err2)
		}
		if tee {
			fmt.Print(string(bytes))
		}
	} else {
		fmt.Print(string(bytes))
	}
}
