package queue

import "testing"

func TestPriorityQueue(t *testing.T) {
	pq := New[int]()

	if pq.Pop() != nil {
		t.Fatal("pop on an empty queue should return nil")
	}

	item0 := &Item[int]{
		Position:   1,
		Step:       nil,
		ErrorsLeft: 1,
	}
	pq.Add(nil) // should be safely ignored¨
	pq.Add(item0)

	if pq.Pop() != item0 {
		t.Fatal("item 0 should of been returned")
	}

	if pq.Pop() != nil {
		t.Fatal("queue should be empty and nil should of been returned")
	}
}
