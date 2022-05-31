package colors

import (
	"github.com/marcadamsge/gofuzzy/fuzzy"
	"github.com/marcadamsge/gofuzzy/trie"
)

import "testing"

func TestTrie(t *testing.T) {
	// Create a new Trie and specify the type that it will be storing
	myTrie := trie.New[string]()

	// Somehow load your dataset
	blue := "blue"
	green := "green"
	black := "black"

	// Define how your data points should be merged together
	// This is used by the trie.Insert function in case there's duplicates
	// Here we simply take one of the non nil values,
	// we know there won't be duplicates anyway
	combineFunction := func(t1 *string, t2 *string) *string {
		if t1 != nil {
			return t1
		}

		return t2
	}

	// Index data in Trie
	myTrie.Insert(blue, &blue, combineFunction)
	myTrie.Insert(green, &green, combineFunction)
	myTrie.Insert(black, &black, combineFunction)

	// Define how the Fuzzy search algorithm should collect data
	// This lets you define:
	//   1. When enough data has been collected
	//   2. How the data should be collected
	myCollector := fuzzy.NewListCollector[string](3)

	// Run the search
	fuzzy.Search[string](myTrie, "bue", 1, myCollector)

	result := myCollector.Results
	if len(result) != 1 || result[0].Value != &blue {
		t.Fatal("unexpected result")
	}
}
