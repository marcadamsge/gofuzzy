package queue

import "container/heap"

type PriorityQueue[T any] struct {
	itemArray *itemArray[T]
}

func New[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		itemArray: &itemArray[T]{},
	}
}

// Add an element to the priority queue. If item is nil, it's ignored.
func (pq *PriorityQueue[T]) Add(item *Item[T]) {
	if item != nil {
		heap.Push(pq.itemArray, item)
	}
}

// Pop and element from the priority queue. If the queue is empty, nil is returned.
func (pq *PriorityQueue[T]) Pop() *Item[T] {
	if pq.itemArray.Len() >= 1 {
		return heap.Pop(pq.itemArray).(*Item[T])
	}

	return nil
}
