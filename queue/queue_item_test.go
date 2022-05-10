package queue

import "testing"

func TestQueueItem(t *testing.T) {
	testQueueItem := &itemArray[int]{}
	if testQueueItem.Len() != 0 {
		t.Fatalf("unexpected length: %d", testQueueItem.Len())
	}

	item0 := &Item[int]{
		Position:   2,
		Step:       nil,
		ErrorsLeft: 1,
	}

	testQueueItem.Push(item0)

	if testQueueItem.Len() != 1 {
		t.Fatalf("unexpected length: %d", testQueueItem.Len())
	}

	item1 := &Item[int]{
		Position:   1,
		Step:       nil,
		ErrorsLeft: 1,
	}

	testQueueItem.Push(item1)

	if testQueueItem.Len() != 2 {
		t.Fatalf("unexpected length: %d", testQueueItem.Len())
	}

	if testQueueItem.Less(0, 1) != true {
		t.Fatal("unexpected result")
	}

	testQueueItem.Swap(0, 1)
	if testQueueItem.Less(0, 1) != false {
		t.Fatal("unexpected result")
	}

	el := testQueueItem.Pop()
	elItem, ok := el.(*Item[int])
	if !ok || elItem != item0 {
		t.Fatal("unexpected value returned")
	}
}

func TestQueueItemLessOp(t *testing.T) {
	testIA := itemArray[int]{
		&Item[int]{
			Position:   2,
			Step:       nil,
			ErrorsLeft: 1,
		},
		&Item[int]{
			Position:   2,
			Step:       nil,
			ErrorsLeft: 1,
		},
	}

	if testIA.Less(0, 1) != false || testIA.Less(1, 0) != false {
		t.Fatal("item 0 and 1 should be equal")
	}

	testIA = itemArray[int]{
		&Item[int]{
			Position:   2,
			Step:       nil,
			ErrorsLeft: 1,
		},
		&Item[int]{
			Position:   1,
			Step:       nil,
			ErrorsLeft: 1,
		},
	}

	if testIA.Less(0, 1) != true || testIA.Less(1, 0) != false {
		t.Fatal("item 0 should be less then item 1")
	}

	testIA = itemArray[int]{
		&Item[int]{
			Position:   2,
			Step:       nil,
			ErrorsLeft: 1,
		},
		&Item[int]{
			Position:   1,
			Step:       nil,
			ErrorsLeft: 2,
		},
	}

	if testIA.Less(0, 1) != false || testIA.Less(1, 0) != true {
		t.Fatal("item 1 should be less then item 0")
	}
}
