package fuzzy

import (
	"context"
	"github.com/marcadamsge/gofuzzy/queue"
	"github.com/marcadamsge/gofuzzy/trie"
	"sync/atomic"
)

// Search a fuzzy match on the trie until collector.Done() is true or there is no more match given the Levenshtein distance.
// Search calls collector.Collect first with the closest match, and then the second closest, etc...
func Search[T any](trie *trie.Trie[T], word string, distance int, collector ResultCollector[T]) {
	doneChannel := make(chan struct{}, 1)
	var continueRun int32 = 1
	continueRunPtr := &continueRun

	search[T](continueRunPtr, doneChannel, trie, word, distance, collector)
}

// SearchAsync does the same as Search but may also stop when the context cancels the search.
// Canceling a search does not delete all the matches already found.
// The output channel is written to when the search is completed to notify that the search is finished.
func SearchAsync[T any](ctx context.Context, trie *trie.Trie[T], word string, distance int, collector ResultCollector[T]) <-chan struct{} {
	cancelChannel := ctx.Done()
	searchDoneChannel := make(chan struct{}, 1)

	var continueRun int32 = 1
	continueRunPtr := &continueRun

	if cancelChannel != nil {
		doneChannel := make(chan struct{}, 1)

		go func() {
			select {
			case <-cancelChannel:
				atomic.StoreInt32(continueRunPtr, 0)
				// wait for the search to cleanly before notifying that the search is done
				doneChannel <- <-searchDoneChannel
				return
			case <-searchDoneChannel:
				doneChannel <- struct{}{}
				return
			}
		}()

		go search[T](continueRunPtr, searchDoneChannel, trie, word, distance, collector)
		return doneChannel
	}

	go search[T](continueRunPtr, searchDoneChannel, trie, word, distance, collector)
	return searchDoneChannel
}

func search[T any](continueRun *int32, doneChannel chan<- struct{}, node *trie.Trie[T], str string, distance int, collector ResultCollector[T]) {
	priorityQueue := queue.New[T]()
	priorityQueue.Add(&queue.Item[T]{
		Position:   0,
		Step:       node,
		ErrorsLeft: distance,
	})

	runes := []rune(str)
	resultSet := make(map[*trie.Trie[T]]struct{})
	maxPosition := len(runes)

	for crtItem := priorityQueue.Pop(); crtItem != nil && atomic.LoadInt32(continueRun) == 1 && !collector.Done(); crtItem = priorityQueue.Pop() {
		if crtItem.ErrorsLeft > 0 && maxPosition > crtItem.Position {
			// a character was randomly changed with another one
			crtItem.Step.Iterate(func(r rune, trie *trie.Trie[T]) {
				if r != runes[crtItem.Position] {
					priorityQueue.Add(&queue.Item[T]{
						Position:   crtItem.Position + 1,
						Step:       trie,
						ErrorsLeft: crtItem.ErrorsLeft - 1,
					})
				}
			})

			// a character was inserted but shouldn't be there
			priorityQueue.Add(&queue.Item[T]{
				Position:   crtItem.Position + 1,
				Step:       crtItem.Step,
				ErrorsLeft: crtItem.ErrorsLeft - 1,
			})
		}

		// a character was removed
		if crtItem.ErrorsLeft > 0 {
			crtItem.Step.Iterate(func(r rune, trie *trie.Trie[T]) {
				priorityQueue.Add(&queue.Item[T]{
					Position:   crtItem.Position,
					Step:       trie,
					ErrorsLeft: crtItem.ErrorsLeft - 1,
				})
			})
		}

		// two adjacent characters were swapped
		if crtItem.ErrorsLeft > 0 && maxPosition-1 > crtItem.Position {
			step1 := crtItem.Step.Step(runes[crtItem.Position+1])
			if step1 != nil {
				step2 := step1.Step(runes[crtItem.Position])
				if step2 != nil {
					priorityQueue.Add(&queue.Item[T]{
						Position:   crtItem.Position + 2,
						Step:       step2,
						ErrorsLeft: crtItem.ErrorsLeft - 1,
					})
				}
			}
		}

		// test if we're in a final state
		if maxPosition == crtItem.Position && crtItem.Step.Value != nil {
			_, resultAlreadyReturned := resultSet[crtItem.Step]

			if !resultAlreadyReturned {
				collector.Collect(crtItem.Step.Value, distance-crtItem.ErrorsLeft)
				resultSet[crtItem.Step] = struct{}{}
			}
		}

		// try stepping out once
		if maxPosition > crtItem.Position {
			nextItem := crtItem.Step.Step(runes[crtItem.Position])
			if nextItem != nil {
				priorityQueue.Add(&queue.Item[T]{
					Position:   crtItem.Position + 1,
					Step:       nextItem,
					ErrorsLeft: crtItem.ErrorsLeft,
				})
			}
		}
	}

	doneChannel <- struct{}{}
}
