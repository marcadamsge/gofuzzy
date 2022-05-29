package fuzzy

type ResultCollector[T any] interface {
	// Collect result when the Fuzzy function founds a match.
	// The Search function calls Collect the first time with the closest match, then the second closest, etc...
	// For example the results of the ListCollector are ordered by distance from least to most.
	Collect(t *T, distance int)
	// Done tells the fuzzy search it can stop
	Done() bool
}

// ListCollector collects all the results in a list until MaxResult number of results are collected
// If MaxResult < 0, then the list collector will collect forever
type ListCollector[T any] struct {
	MaxResult int
	Results   []Result[T]
}

type Result[T any] struct {
	Value    *T
	Distance int
}

func NewListCollector[T any](maxResult int) *ListCollector[T] {
	return &ListCollector[T]{
		MaxResult: maxResult,
		Results:   make([]Result[T], 0, 0),
	}
}

func (lc *ListCollector[T]) Collect(t *T, distance int) {
	if t != nil {
		lc.Results = append(lc.Results, Result[T]{
			Value:    t,
			Distance: distance,
		})
	}
}

func (lc *ListCollector[T]) Done() bool {
	if lc.MaxResult >= 0 {
		return len(lc.Results) >= lc.MaxResult
	}

	return false
}

func NewCountCollector[T any](maxResult int) *CountCollector[T] {
	return &CountCollector[T]{
		MaxResult:   maxResult,
		ResultCount: 0,
	}
}

// CountCollector simply counts the number of results until MaxResult number is collected.
// The actual results are discarded.
type CountCollector[T any] struct {
	MaxResult   int
	ResultCount int
}

func (cc *CountCollector[T]) Collect(t *T, distance int) {
	if t != nil {
		cc.ResultCount++
	}
}

func (cc *CountCollector[T]) Done() bool {
	return cc.ResultCount >= cc.MaxResult
}
