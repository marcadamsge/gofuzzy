package queue

import "github.com/marcadamsge/gofuzzy/trie"

type Item[T any] struct {
	// our position in the input string
	Position int
	// current step in the Trie we are exploring
	Step *trie.Trie[T]
	// number of errors that can still be made
	ErrorsLeft int
}

type itemArray[T any] []*Item[T]

// we implement the heap.Interface methods

func (ia itemArray[T]) Len() int {
	return len(ia)
}

func (ia itemArray[T]) Less(i, j int) bool {
	pqi := ia[i]
	pqj := ia[j]

	if pqi.ErrorsLeft != pqj.ErrorsLeft {
		return pqi.ErrorsLeft > pqj.ErrorsLeft
	}

	return pqi.Position > pqj.Position
}

func (ia itemArray[T]) Swap(i, j int) {
	ia[i], ia[j] = ia[j], ia[i]
}

func (ia *itemArray[T]) Push(x any) {
	item := x.(*Item[T])
	*ia = append(*ia, item)
}

func (ia *itemArray[T]) Pop() any {
	old := *ia
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*ia = old[0 : n-1]
	return item
}
