package trie

type Trie[T any] struct {
	children map[rune]*Trie[T]
	Value    *T
}

func New[T any]() *Trie[T] {
	return &Trie[T]{
		children: make(map[rune]*Trie[T]),
		Value:    nil,
	}
}

// Insert a string into the trie. str may be empty, in which case the value of the current trie is updated.
// The combineValues function merges the value of the trie being updated with the new value, this is a way
// to avoid overwriting previously inserted values.
// This function is pretty basic, if you need more control over how the values are inserted
// (like storing prefixes as well for example), it's better to use StepOrCreate directly.
func (trie *Trie[T]) Insert(str string, value *T, combineValues func(t1 *T, t2 *T) *T) {
	crtTrie := trie
	for _, r := range []rune(str) {
		crtTrie = crtTrie.StepOrCreate(r)
	}
	crtTrie.Value = combineValues(crtTrie.Value, value)
}

// Step out with the rune r and return the next Trie or nil if it does not exist.
func (trie *Trie[T]) Step(r rune) *Trie[T] {
	return trie.children[r]
}

// StepOrCreate the next Trie if it does not exist.
func (trie *Trie[T]) StepOrCreate(r rune) *Trie[T] {
	step, ok := trie.children[r]
	if ok {
		return step
	}

	out := &Trie[T]{
		children: make(map[rune]*Trie[T]),
		Value:    nil,
	}
	trie.children[r] = out
	return out
}

// Iterate over all the children of this trie.
func (trie *Trie[T]) Iterate(iterationFunction func(r rune, trie *Trie[T])) {
	for r, tr := range trie.children {
		iterationFunction(r, tr)
	}
}
