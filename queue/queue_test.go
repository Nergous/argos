package queue

import "testing"

func assertNewEmpty[T any](t *testing.T) {
	t.Helper()

	q := New[T]()
	if got := q.Len(); got != 0 {
		t.Errorf("Len() = %d, want 0", got)
	}
	if got := cap(q.data); got != 0 {
		t.Errorf("cap(data) = %d, want 0", got)
	}
}

func TestQueue_New(t *testing.T) {
	t.Run("int", assertNewEmpty[int])
	t.Run("string", assertNewEmpty[string])
	t.Run("float32", assertNewEmpty[float32])
}

func TestQueue_NewWithCap(t *testing.T) {
	tests := []struct {
		name string
		cap  int
	}{
		{name: "zero cap", cap: 0},
		{name: "one", cap: 1},
		{name: "ten", cap: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewWithCap[int](tt.cap)
			if got := q.Len(); got != 0 {
				t.Errorf("Len() = %d, want 0", got)
			}
			if got := cap(q.data); got != tt.cap {
				t.Errorf("cap(data) = %d, want %d", got, tt.cap)
			}
		})
	}
}

func TestQueue_PushPop(t *testing.T) {
	tests := []struct {
		name string
		push []int
		want []int
	}{
		{name: "empty", push: nil, want: nil},
		{name: "single element", push: []int{42}, want: []int{42}},
		{name: "fifo order", push: []int{1, 2, 3}, want: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New[int]()
			for _, v := range tt.push {
				q.Push(v)
			}

			for i, want := range tt.want {
				got, ok := q.Pop()
				if !ok {
					t.Fatalf("Pop() #%d: ok = false, want true", i)
				}
				if got != want {
					t.Errorf("Pop() #%d = %d, want %d", i, got, want)
				}
			}

			if got, ok := q.Pop(); ok {
				t.Errorf("Pop() on empty: got %d, ok = true, want ok = false", got)
			}
		})
	}
}

func TestQueue_PushPopInterleaved(t *testing.T) {
	// FIFO order must hold even when pushes and pops are interleaved.
	q := New[int]()

	q.Push(1)
	q.Push(2)

	if got, ok := q.Pop(); !ok || got != 1 {
		t.Fatalf("Pop() = (%d, %v), want (1, true)", got, ok)
	}

	q.Push(3)

	for _, want := range []int{2, 3} {
		got, ok := q.Pop()
		if !ok {
			t.Fatalf("Pop(): ok = false, want true")
		}
		if got != want {
			t.Errorf("Pop() = %d, want %d", got, want)
		}
	}

	if got, ok := q.Pop(); ok {
		t.Errorf("Pop() on empty: got %d, ok = true, want ok = false", got)
	}
}

func TestQueue_PopReleasesReference(t *testing.T) {
	// Pop must zero the vacated slot so the backing array no longer retains a
	// reference to the popped element (otherwise it leaks until reallocation).
	q := New[*int]()
	a := 1
	q.Push(&a)

	// alias the same backing array so we can inspect the slot after Pop slices it off.
	backing := q.data[:1]

	if _, ok := q.Pop(); !ok {
		t.Fatalf("Pop() ok = false, want true")
	}
	if backing[0] != nil {
		t.Errorf("backing slot 0 = %p, want nil (Pop must release the reference)", backing[0])
	}
}

func TestQueue_Peek(t *testing.T) {
	t.Run("empty returns zero value and false", func(t *testing.T) {
		q := New[int]()
		got, ok := q.Peek()
		if ok {
			t.Errorf("Peek() ok = true, want false")
		}
		if got != 0 {
			t.Errorf("Peek() = %d, want 0 (zero value)", got)
		}
	})

	t.Run("returns front without removing it", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		q.Push(2)

		got, ok := q.Peek()
		if !ok {
			t.Fatalf("Peek() ok = false, want true")
		}
		if got != 1 {
			t.Errorf("Peek() = %d, want 1 (front of queue)", got)
		}
		if got := q.Len(); got != 2 {
			t.Errorf("Len() after Peek() = %d, want 2 (Peek must not modify the queue)", got)
		}
	})
}

func TestQueue_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		pushN int
		want  bool
	}{
		{name: "new queue is empty", pushN: 0, want: true},
		{name: "after push is not empty", pushN: 1, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New[int]()
			for i := 0; i < tt.pushN; i++ {
				q.Push(i)
			}
			if got := q.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue_Cap(t *testing.T) {
	t.Run("new queue has zero cap", func(t *testing.T) {
		q := New[int]()
		if got := q.Cap(); got != 0 {
			t.Errorf("Cap() = %d, want 0", got)
		}
	})

	t.Run("reflects preallocated capacity", func(t *testing.T) {
		q := NewWithCap[int](8)
		if got := q.Cap(); got != 8 {
			t.Errorf("Cap() = %d, want 8", got)
		}
		if got := q.Len(); got != 0 {
			t.Errorf("Len() = %d, want 0", got)
		}
	})

	t.Run("never less than len while growing", func(t *testing.T) {
		q := New[int]()
		for i := range 100 {
			q.Push(i)
			if c, l := q.Cap(), q.Len(); c < l {
				t.Fatalf("after %d pushes: Cap() = %d < Len() = %d", i+1, c, l)
			}
		}
	})

	t.Run("no realloc while pushing within preallocated cap", func(t *testing.T) {
		q := NewWithCap[int](4)
		for i := range 4 {
			q.Push(i)
		}
		if got := q.Cap(); got != 4 {
			t.Errorf("Cap() = %d, want 4 (no realloc expected within cap)", got)
		}
	})
}

func TestQueue_Clear(t *testing.T) {
	q := NewWithCap[int](8)
	for i := range 5 {
		q.Push(i)
	}

	q.Clear()

	if got := q.Len(); got != 0 {
		t.Errorf("Len() after Clear() = %d, want 0", got)
	}
	if !q.IsEmpty() {
		t.Errorf("IsEmpty() after Clear() = false, want true")
	}
	if got := q.Cap(); got != 0 {
		t.Errorf("Cap() after Clear() = %d, want 0 (backing array released)", got)
	}
	if v, ok := q.Pop(); ok {
		t.Errorf("Pop() after Clear() = (%d, true), want ok = false", v)
	}

	// queue stays usable after Clear
	q.Push(99)
	if v, ok := q.Pop(); !ok || v != 99 {
		t.Errorf("Pop() after reuse = (%d, %v), want (99, true)", v, ok)
	}
}

func TestQueue_Reset(t *testing.T) {
	t.Run("empties queue but preserves capacity", func(t *testing.T) {
		q := New[int]()
		for i := range 10 {
			q.Push(i)
		}
		capBefore := q.Cap()

		q.Reset()

		if got := q.Len(); got != 0 {
			t.Errorf("Len() after Reset() = %d, want 0", got)
		}
		if !q.IsEmpty() {
			t.Errorf("IsEmpty() after Reset() = false, want true")
		}
		if got := q.Cap(); got != capBefore {
			t.Errorf("Cap() after Reset() = %d, want %d (capacity must be preserved)", got, capBefore)
		}
	})

	t.Run("zeroes backing array to release references", func(t *testing.T) {
		q := New[*int]()
		a, b, c := 1, 2, 3
		q.Push(&a)
		q.Push(&b)
		q.Push(&c)

		q.Reset()

		// white-box: the array is retained, so inspect every slot up to cap.
		full := q.data[:cap(q.data)]
		for i, p := range full {
			if p != nil {
				t.Errorf("backing slot %d = %p, want nil (Reset must zero elements)", i, p)
			}
		}
	})

	t.Run("queue stays usable after Reset", func(t *testing.T) {
		q := New[int]()
		q.Push(1)
		q.Reset()
		q.Push(42)

		if v, ok := q.Pop(); !ok || v != 42 {
			t.Errorf("Pop() after reuse = (%d, %v), want (42, true)", v, ok)
		}
		if v, ok := q.Pop(); ok {
			t.Errorf("Pop() on emptied queue = (%d, true), want ok = false", v)
		}
	})
}
