package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/marcadamsge/gofuzzy/fuzzy"
	"github.com/marcadamsge/gofuzzy/gen"
	"github.com/marcadamsge/gofuzzy/trie"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func fuzzySearchPerfTest(
	geoNamesReader io.Reader,
	genNamesTrie *trie.Trie[Entry],
	seed int64,
	threads int,
	numberOfLines uint32,
	maxResult int,
) error {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	println("starting performance test...")

	if threads < 1 {
		return errors.New("the number of threads should be at least 1")
	}

	if maxResult < 1 {
		return errors.New("max results should be at least 1")
	}

	var workerWaitGroup sync.WaitGroup

	workerWaitGroup.Add(threads + 1)

	readChannel := make(chan string, 4*threads)
	perfResultChannel := make(chan testOutput, 4*threads)

	// perfTestReader will close the readChannel once it finished reading
	go perfTestReader(readChannel, geoNamesReader, numberOfLines, &workerWaitGroup)

	randGen := rand.New(rand.NewSource(seed))

	for i := 0; i < threads; i++ {
		go perfTestWorker(
			readChannel,
			perfResultChannel,
			genNamesTrie,
			// random generator is not thread safe, so we create one per worker
			rand.New(rand.NewSource(randGen.Int63())),
			alphabet,
			maxResult,
			&workerWaitGroup,
		)
	}

	var collectorWaitGroup sync.WaitGroup
	collectorWaitGroup.Add(1)
	go perfTestCollectResults(perfResultChannel, &collectorWaitGroup)

	// perfTestWorker do not close the perfResultChannel, so we do it here
	workerWaitGroup.Wait()
	close(perfResultChannel)

	collectorWaitGroup.Wait()
	return nil
}

func perfTestReader(
	outputChannel chan<- string,
	geoNamesReader io.Reader,
	numberOfLines uint32,
	waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()
	defer close(outputChannel)

	geoNamesScanner := bufio.NewScanner(geoNamesReader)
	linesParsed := uint32(0)
	lastPercent := uint32(0)

	print("0% done")

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

		outputChannel <- name
		linesParsed++

		percentDone := linesParsed * 100 / numberOfLines
		if percentDone != lastPercent {
			lastPercent = percentDone
			fmt.Printf("\033[1K\r%d%% done", percentDone)
		}
	}
	println()

	if geoNamesScanner.Err() != nil {
		fmt.Printf("got error while reading geonames input file: %s\n", geoNamesScanner.Err())
	}
}

type testOutput struct {
	fuzzyLength  int
	timeTaken    time.Duration
	resultsFound int
}

func perfTestWorker(
	inputChannel <-chan string,
	outputChannel chan<- testOutput,
	genNamesTrie *trie.Trie[Entry],
	randGen gen.RandIntGenerator,
	alphabet []rune,
	maxResults int,
	waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()

	for name, ok := <-inputChannel; ok; name, ok = <-inputChannel {
		var maxDistance int
		if len(name) <= 2 {
			maxDistance = 0
		} else if len(name) <= 5 {
			maxDistance = 1
		} else {
			maxDistance = 2
		}

		fuzzyName := gen.RandomFuzzyErrors(name, randGen, maxDistance, alphabet)
		collector := fuzzy.NewCountCollector[Entry](maxResults)

		start := time.Now()
		fuzzy.Search[Entry](context.Background(), genNamesTrie, fuzzyName, maxDistance, collector)
		end := time.Now()

		outputChannel <- testOutput{
			fuzzyLength:  len(fuzzyName),
			timeTaken:    end.Sub(start),
			resultsFound: collector.ResultCount,
		}
	}
}

func perfTestCollectResults(inputChannel <-chan testOutput, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	result, ok := <-inputChannel

	if !ok {
		return
	}

	fuzzyLengthMin := result.fuzzyLength
	fuzzyLengthMax := result.fuzzyLength
	fuzzyLengthTotal := uint64(result.fuzzyLength)

	timeMin := result.timeTaken
	timeMax := result.timeTaken
	timeTotal := uint64(result.timeTaken)

	count := uint32(1)

	for ; ok; result, ok = <-inputChannel {
		fuzzyLengthMin = minInt(fuzzyLengthMin, result.fuzzyLength)
		fuzzyLengthMax = maxInt(fuzzyLengthMax, result.fuzzyLength)
		fuzzyLengthTotal = fuzzyLengthTotal + uint64(result.fuzzyLength)

		timeMin = minInt(timeMin, result.timeTaken)
		timeMax = maxInt(timeMax, result.timeTaken)
		timeTotal = timeTotal + uint64(result.timeTaken.Nanoseconds())

		count++
	}

	averageLength := float64(fuzzyLengthTotal) / float64(count)
	averageTime := float64(timeTotal) / float64(count)

	println("results:")
	fmt.Printf("length: %d min, %d max, %f average\n", fuzzyLengthMin, fuzzyLengthMax, averageLength)
	fmt.Printf("time (in nano second): %d min, %d max, %f average\n", timeMin.Nanoseconds(), timeMax.Nanoseconds(), averageTime)
}

func maxInt[T int | time.Duration](a T, b T) T {
	if a >= b {
		return a
	}

	return b
}

func minInt[T int | time.Duration](a T, b T) T {
	if a <= b {
		return a
	}

	return b
}
