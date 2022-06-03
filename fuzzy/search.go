package fuzzy

import (
	"context"
	"github.com/marcadamsge/gofuzzy/queue"
	"github.com/marcadamsge/gofuzzy/trie"
	"sync/atomic"
)

// Search a fuzzy match on the trie until collector.Done() is true or there is no more match given the Levenshtein distance.
// Search calls collector.Collect first with the closest match, and then the second closest, etc...
func Search[T any](ctx context.Context, trie *trie.Trie[T], word string, distance int, collector ResultCollector[T]) {
	var continueRun int32 = 1
	continueRunPtr := &continueRun

	cancelChannel := ctx.Done()
	if cancelChannel != nil {
		searchDoneChannel := make(chan struct{}, 1)
		defer close(searchDoneChannel)

		go func() {
			select {
			case <-cancelChannel:
				atomic.StoreInt32(continueRunPtr, 0)
				return
			case <-searchDoneChannel:
				return
			}
		}()
	}

	search[T](continueRunPtr, trie, word, distance, collector)
}

func search[T any](continueRun *int32, node *trie.Trie[T], str string, distance int, collector ResultCollector[T]) {
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
}
