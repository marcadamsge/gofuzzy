package fuzzy

import (
	"github.com/marcadamsge/gofuzzy/trie"
	"reflect"
	"runtime/debug"
	"testing"
)

func TestSearch(t *testing.T) {
	testTrie := trie.New[string]()
	word1 := "cat"
	word2 := "tat"
	word3 := "dog"

	combineFunction := func(t1 *string, t2 *string) *string {
		if t1 != nil {
			return t1
		}

		return t2
	}

	testTrie.Insert(word1, &word1, combineFunction)
	testTrie.Insert(word2, &word2, combineFunction)
	testTrie.Insert(word3, &word3, combineFunction)

	checkResult(
		t,
		testTrie, "cat", 0, 1,
		[]Result[string]{
			{
				Value:    &word1,
				Distance: 0,
			},
		},
	)

	checkResult(
		t,
		testTrie, "og", 1, 1,
		[]Result[string]{
			{
				Value:    &word3,
				Distance: 1,
			},
		},
	)

	checkResult(
		t,
		testTrie, "do", 1, 1,
		[]Result[string]{
			{
				Value:    &word3,
				Distance: 1,
			},
		},
	)

	checkResult(
		t,
		testTrie, "dgo", 1, 1,
		[]Result[string]{
			{
				Value:    &word3,
				Distance: 1,
			},
		},
	)

	checkResult(
		t,
		testTrie, "dogg", 1, 1,
		[]Result[string]{
			{
				Value:    &word3,
				Distance: 1,
			},
		},
	)

	checkResult(
		t,
		testTrie, "dogd", 1, 1,
		[]Result[string]{
			{
				Value:    &word3,
				Distance: 1,
			},
		},
	)

	checkResult(
		t,
		testTrie, "dod", 1, 1,
		[]Result[string]{
			{
				Value:    &word3,
				Distance: 1,
			},
		},
	)

	checkResult(
		t,
		testTrie, "cat", 1, 2,
		[]Result[string]{
			{
				Value:    &word1,
				Distance: 0,
			},
			{
				Value:    &word2,
				Distance: 1,
			},
		},
	)

	checkResult(
		t,
		testTrie, "cat", 1, 1,
		[]Result[string]{
			{
				Value:    &word1,
				Distance: 0,
			},
		},
	)

	checkResult(
		t,
		testTrie, "cat", 3, 4,
		[]Result[string]{
			{
				Value:    &word1,
				Distance: 0,
			},
			{
				Value:    &word2,
				Distance: 1,
			},
			{
				Value:    &word3,
				Distance: 3,
			},
		},
	)
}

func checkResult(t *testing.T, trie *trie.Trie[string], word string, distance int, maxResults int, expectedResult []Result[string]) {
	collector := NewListCollector[string](maxResults)
	Search[string](trie, word, distance, collector)

	if !reflect.DeepEqual(collector.Results, expectedResult) {
		t.Log(string(debug.Stack()))
		t.Fatal("unexpected result")
	}
}
