package colors

import (
	"github.com/marcadamsge/gofuzzy/fuzzy"
	"github.com/marcadamsge/gofuzzy/trie"
)

import "testing"

func TestTrie(t *testing.T) {
	myTrie := trie.New[string]()

	blue := "blue"
	green := "green"
	black := "black"

	combineFunction := func(t1 *string, t2 *string) *string {
		if t1 != nil {
			return t1
		}

		return t2
	}

	myTrie.Insert(blue, &blue, combineFunction)
	myTrie.Insert(green, &green, combineFunction)
	myTrie.Insert(black, &black, combineFunction)

	myCollector := fuzzy.NewListCollector[string](3)
	fuzzy.Search[string](myTrie, "bue", 1, myCollector)

	result := myCollector.Results

	if len(result) != 1 || result[0].Value != &blue {
		t.Fatal("unexpected result")
	}
}
