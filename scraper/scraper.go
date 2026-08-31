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
		debug           = false
		noLog           = false
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
	if slices.Contains(os.Args, "--no-log") {
		noLog = true
	}
	if slices.Contains(os.Args, "--debug") {
		debug = true
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
	}
	if noLog {
		log.Disable()
	}
	if debug {
		log.EnableDebug()
		log.Debugln("debug logging enabled")
	}

	Execute(useMock, pretty, outFile, tee, truncateItemsTo)
}

func Execute(useMock, pretty bool, outFile string, tee bool, truncateItemsTo int) {
	if useMock {
		log.Println("Using mock fetcher")
	}
	if pretty {
		log.Println("Pretty printing")
	}
	if outFile != "" {
		log.Printf("Using output file %s\n", outFile)
		if tee {
			log.Println("Also writing output to stdout")
		}
	}
	if truncateItemsTo >= 0 {
		log.Printf("Truncating items to %d\n", truncateItemsTo)
	}

	startTime := time.Now()
	f := fetcher.New(useMock, time.Second)

	log.Debugln("Fetching Collection log table")
	clogTable, err := f.FetchPageSection("Collection_log%2FTable", "Table")
	if err != nil {
		panic(err)
	}

	log.Debugln("Parsing Collection log table")
	items, err := parser.ParseCollectionLogTable(clogTable)
	if err != nil {
		panic(err)
	}
	log.Debugf("Parsed %d Collection log items\n", len(items))
	if truncateItemsTo >= 0 {
		items = items[:truncateItemsTo]
		log.Debugf("Truncating items to %d\n", truncateItemsTo)
	}

	resultSet := results.New()
	for _, item := range items {
		log.Debugf("Fetching %s item sources\n", item.Name)
		sec, err2 := f.FetchPageSection(item.Link, "Item sources")
		if err2 != nil {
			if strings.HasSuffix(err2.Error(), "section not found") {
				log.Debugf("Found no %s item sources\n", item.Name)
				resultSet.Add(parser.ItemSourceNone(item))
				continue
			} else {
				panic(err2)
			}
		}
		sources, err2 := parser.ParseItemSources(item, sec)
		if err2 != nil {
			panic(err2)
		}
		log.Debugf("Found %d %s item sources: %s\n", len(sources), item.Name, sources)
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
