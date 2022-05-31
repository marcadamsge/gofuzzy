package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	geoNamesFileName := flag.String("geo", "", "GeoNames file to parse")
	threads := flag.Int("threads", runtime.NumCPU(), "number of threads to use for the test")
	maxResults := flag.Int("n", 1, "max number of results per test")
	flag.Parse()

	if geoNamesFileName == nil || *geoNamesFileName == "" {
		println("-geo flag is required")
		os.Exit(1)
	}

	if threads == nil || *threads < 1 {
		println("at least one thread is required")
		os.Exit(1)
	}

	if maxResults == nil || *maxResults < 1 {
		println("at least on result is needed")
		os.Exit(1)
	}

	geoNamesReader, err := os.Open(*geoNamesFileName)
	if err != nil {
		fmt.Printf("failed to open geonames file with error: %s\n", err.Error())
		os.Exit(1)
	}
	defer geoNamesReader.Close()

	geoNamesTrie, numberOfLines, err := parseGeoNamesFile(geoNamesReader)
	if err != nil {
		fmt.Printf("failed to read geonames file with error: %s\n", err.Error())
		os.Exit(1)
	}

	triggerGC()

	_, err = geoNamesReader.Seek(0, 0)
	if err != nil {
		fmt.Printf("failed to seek at the beginning of the geonames file with error: %s\n", err.Error())
		os.Exit(1)
	}

	err = fuzzySearchPerfTest(
		geoNamesReader,
		geoNamesTrie,
		time.Now().UnixNano(),
		*threads,
		numberOfLines,
		*maxResults,
	)
	if err != nil {
		fmt.Printf("failed to read geonames file with error: %s\n", err.Error())
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("got error while parsing the geonames file: %s\n", err.Error())
		os.Exit(1)
	}
}

func triggerGC() {
	println("triggering manual GC...")
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated Memory = %v MiB\n", m.Alloc/1024/1024)
}
