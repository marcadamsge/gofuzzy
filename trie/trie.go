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
	runes := []rune(str)
	crtTrie := trie
	if len(runes) >= 2 {
		for _, r := range runes[:len(runes)-1] {
			crtTrie = crtTrie.stepOrCreate(r)
		}
	}

	if len(runes) >= 1 {
		crtTrie.createOrMerge(runes[len(runes)-1], value, combineValues)
	}
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

func (trie *Trie[T]) createOrMerge(r rune, value *T, combineValues func(t1 *T, t2 *T) *T) {
	step, ok := trie.children[r]
	if ok {
		step.Value = combineValues(step.Value, value)
	} else {
		out := &Trie[T]{
			children: make(map[rune]*Trie[T]),
			Value:    value,
		}
		trie.children[r] = out
	}
}
