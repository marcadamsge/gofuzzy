package fuzzy

import (
	"context"
	"github.com/marcadamsge/gofuzzy/queue"
	"github.com/marcadamsge/gofuzzy/trie"
	"sync/atomic"
)

type Result[T any] struct {
	Value    *T
	Distance int
}

func SearchAsync[T any](ctx context.Context, trie *trie.Trie[T], word string, distance int, maxResult int) <-chan []Result[T] {
	cancelChannel := ctx.Done()
	resultChannel := make(chan []Result[T], 1)

	var continueRun int32 = 1
	continueRunPtr := &continueRun

	if cancelChannel != nil {
		outputChannel := make(chan []Result[T], 1)

		go func() {
			select {
			case <-cancelChannel:
				atomic.StoreInt32(continueRunPtr, 0)
				outputChannel <- <-resultChannel
				return
			case output := <-resultChannel:
				outputChannel <- output
				return
			}
		}()

		return outputChannel
	}

	go search[T](continueRunPtr, resultChannel, trie, word, distance, maxResult)

	return resultChannel
}

func Search[T any](trie *trie.Trie[T], word string, distance int, maxResult int) []Result[T] {
	resultChannel := make(chan []Result[T], 1)
	var continueRun int32 = 1
	continueRunPtr := &continueRun

	search[T](continueRunPtr, resultChannel, trie, word, distance, maxResult)

	return <-resultChannel
}

func search[T any](continueRun *int32, resultChannel chan<- []Result[T], node *trie.Trie[T], str string, distance int, maxResult int) {
	priorityQueue := queue.New[T]()
	priorityQueue.Add(&queue.Item[T]{
		Position:   0,
		Step:       node,
		ErrorsLeft: distance,
	})

	runes := []rune(str)
	output := make([]Result[T], 0, 0)
	resultSet := make(map[*trie.Trie[T]]interface{})
	maxPosition := len(runes)

	for crtItem := priorityQueue.Pop(); crtItem != nil && atomic.LoadInt32(continueRun) == 1 && len(output) < maxResult; crtItem = priorityQueue.Pop() {
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
				output = append(output, Result[T]{
					Value:    crtItem.Step.Value,
					Distance: distance - crtItem.ErrorsLeft,
				})

				resultSet[crtItem.Step] = nil
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

	resultChannel <- output
}
