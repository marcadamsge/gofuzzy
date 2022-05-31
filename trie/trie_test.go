package trie

import "testing"

func TestTrie(t *testing.T) {
	testTrie := New[int]()
	testCombineFunction := func(i1 *int, i2 *int) *int {
		if i1 != nil && i2 != nil {
			res := *i1 + *i2
			return &res
		}

		if i1 != nil {
			return i1
		}

		return i2
	}

	if testTrie.Value != nil || len(testTrie.children) != 0 {
		t.Fatal("unexpected init state")
	}

	i1 := 0
	testTrie.Insert("", &i1, testCombineFunction)
	if testTrie.Value != &i1 || len(testTrie.children) != 0 {
		t.Fatal("unexpected state")
	}

	i2 := 1
	testTrie.Insert("a", &i2, testCombineFunction)
	if stepOut := testTrie.Step('a'); stepOut == nil || *stepOut.Value != 1 {
		t.Fatal("string 'a' was not added properly")
	}

	if impossibleStep := testTrie.Step('b'); impossibleStep != nil {
		t.Fatal("this step should not be possible")
	}

	i3 := 3
	testTrie.Insert("abc", &i3, testCombineFunction)
	aStep := testTrie.Step('a')
	bStep := aStep.Step('b')
	cStep := bStep.Step('c')

	if *aStep.Value != 1 || bStep.Value != nil || *cStep.Value != 3 {
		t.Fatal("trie was not properly updated for string 'abc'")
	}

	i4 := 3
	testTrie.Insert("abc", &i4, testCombineFunction)
	aStep = testTrie.Step('a')
	bStep = aStep.Step('b')
	cStep = bStep.Step('c')

	if *aStep.Value != 1 || bStep.Value != nil || *cStep.Value != 6 {
		t.Fatal("trie was not properly updated for the second string 'abc'")
	}

	i5 := 4
	testTrie.Insert("⌘", &i5, testCombineFunction)
	if utf8Step := testTrie.Step('⌘'); utf8Step == nil {
		t.Fatal("utf8 character was not properly added")
	}
}
