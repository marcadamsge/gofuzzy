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

func (trie *Trie[T]) Insert(str string, value *T, combineValues func(t1 *T, t2 *T) *T) {
	crtTrie := trie
	for _, r := range []rune(str) {
		crtTrie = crtTrie.stepOrCreate(r)
	}
	crtTrie.Value = combineValues(crtTrie.Value, value)
}

// Step out with the rune r and return the next Trie or nil if it does not exist.
func (trie *Trie[T]) Step(r rune) *Trie[T] {
	return trie.children[r]
}

func (trie *Trie[T]) Iterate(iterationFunction func(r rune, trie *Trie[T])) {
	for r, tr := range trie.children {
		iterationFunction(r, tr)
	}
}

func (trie *Trie[T]) stepOrCreate(r rune) *Trie[T] {
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
