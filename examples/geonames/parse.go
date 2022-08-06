package main

import (
	"bufio"
	"fmt"
	"github.com/marcadamsge/gofuzzy/trie"
	"io"
	"strconv"
	"strings"
	"time"
)

type GeoLocation struct {
	Latitude  float32
	Longitude float32
	Country   string
}

type Entry struct {
	Name string
	// could be two locations have the same name
	LocationSet map[*GeoLocation]struct{}
}

func parseGeoNamesFile(geoNamesReader io.Reader) (*trie.Trie[Entry], uint32, error) {
	geoNamesScanner := bufio.NewScanner(geoNamesReader)
	genNamesTrie := trie.New[Entry]()
	linesParsed := uint32(0)
	startTime := time.Now()

	println("loading file...")

	for geoNamesScanner.Scan() {
		line := strings.Split(geoNamesScanner.Text(), "\t")
		if len(line) < 10 {
			continue
		}

		if line[6] != "P" {
			// we only load cities
			continue
		}

		name := line[1]
		if len(name) == 0 {
			continue
		}

		latitude, err := strconv.ParseFloat(line[4], 32)
		if err != nil {
			return nil, 0, err
		}
		longitude, err := strconv.ParseFloat(line[5], 32)
		if err != nil {
			return nil, 0, err
		}

		location := &GeoLocation{
			Latitude:  float32(latitude),
			Longitude: float32(longitude),
			Country:   line[8],
		}

		linesParsed++

		entry := &Entry{
			Name:        name,
			LocationSet: map[*GeoLocation]struct{}{location: {}},
		}

		genNamesTrie.Insert(name, entry, combineEntries)
	}

	if geoNamesScanner.Err() != nil {
		return nil, 0, geoNamesScanner.Err()
	}

	totalTime := time.Now().Sub(startTime)
	fmt.Printf("dataset loaded in %f seconds\n", totalTime.Seconds())
	fmt.Printf("%d lines parsed, %d elements inserted in the trie\n", linesParsed, linesParsed)

	return genNamesTrie, linesParsed, nil
}

func combineEntries(e1 *Entry, e2 *Entry) *Entry {
	if e1 != nil && e2 != nil {
		for k := range e2.LocationSet {
			e1.LocationSet[k] = struct{}{}
		}

		return e1
	}

	if e1 != nil {
		return e1
	}

	return e2
}
