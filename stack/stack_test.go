package stack

import "testing"

func assertNewEmpty[T any](t *testing.T) {
	t.Helper()

	s := New[T]()
	if got := s.Len(); got != 0 {
		t.Errorf("Len() = %d, want 0", got)
	}
	if got := cap(s.data); got != 0 {
		t.Errorf("cap(data) = %d, want 0", got)
	}
}

func TestStack_New(t *testing.T) {
	t.Run("int", assertNewEmpty[int])
	t.Run("string", assertNewEmpty[string])
	t.Run("float32", assertNewEmpty[float32])
}

func TestStack_NewWithCap(t *testing.T) {
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
			s := NewWithCap[int](tt.cap)
			if got := s.Len(); got != 0 {
				t.Errorf("Len() = %d, want 0", got)
			}
			if got := cap(s.data); got != tt.cap {
				t.Errorf("cap(data) = %d, want %d", got, tt.cap)
			}
		})
	}
}

func TestStack_PushPop(t *testing.T) {
	tests := []struct {
		name string
		push []int
		want []int
	}{
		{name: "empty", push: nil, want: nil},
		{name: "single element", push: []int{42}, want: []int{42}},
		{name: "lifo order", push: []int{1, 2, 3}, want: []int{3, 2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int]()
			for _, v := range tt.push {
				s.Push(v)
			}

			for i, want := range tt.want {
				got, ok := s.Pop()
				if !ok {
					t.Fatalf("Pop() #%d: ok = false, want true", i)
				}
				if got != want {
					t.Errorf("Pop() #%d = %d, want %d", i, got, want)
				}
			}

			if got, ok := s.Pop(); ok {
				t.Errorf("Pop() on empty: got %d, ok = true, want ok = false", got)
			}
		})
	}
}

func TestStack_Peek(t *testing.T) {
	t.Run("empty returns zero value and false", func(t *testing.T) {
		s := New[int]()
		got, ok := s.Peek()
		if ok {
			t.Errorf("Peek() ok = true, want false")
		}
		if got != 0 {
			t.Errorf("Peek() = %d, want 0 (zero value)", got)
		}
	})

	t.Run("returns top without removing it", func(t *testing.T) {
		s := New[int]()
		s.Push(1)
		s.Push(2)

		got, ok := s.Peek()
		if !ok {
			t.Fatalf("Peek() ok = false, want true")
		}
		if got != 2 {
			t.Errorf("Peek() = %d, want 2", got)
		}
		if got := s.Len(); got != 2 {
			t.Errorf("Len() after Peek() = %d, want 2 (Peek must not modify the stack)", got)
		}
	})
}

func TestStack_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		pushN int
		want  bool
	}{
		{name: "new stack is empty", pushN: 0, want: true},
		{name: "after push is not empty", pushN: 1, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int]()
			for i := 0; i < tt.pushN; i++ {
				s.Push(i)
			}
			if got := s.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Cap(t *testing.T) {
	t.Run("new stack has zero cap", func(t *testing.T) {
		s := New[int]()
		if got := s.Cap(); got != 0 {
			t.Errorf("Cap() = %d, want 0", got)
		}
	})

	t.Run("reflects preallocated capacity", func(t *testing.T) {
		s := NewWithCap[int](8)
		if got := s.Cap(); got != 8 {
			t.Errorf("Cap() = %d, want 8", got)
		}
		if got := s.Len(); got != 0 {
			t.Errorf("Len() = %d, want 0", got)
		}
	})

	t.Run("never less than len while growing", func(t *testing.T) {
		s := New[int]()
		for i := range 100 {
			s.Push(i)
			if c, l := s.Cap(), s.Len(); c < l {
				t.Fatalf("after %d pushes: Cap() = %d < Len() = %d", i+1, c, l)
			}
		}
	})

	t.Run("no realloc while pushing within preallocated cap", func(t *testing.T) {
		s := NewWithCap[int](4)
		for i := range 4 {
			s.Push(i)
		}
		if got := s.Cap(); got != 4 {
			t.Errorf("Cap() = %d, want 4 (no realloc expected within cap)", got)
		}
	})
}

func TestStack_All(t *testing.T) {
	collect := func(s *Stack[int]) []int {
		var out []int
		for v := range s.All() {
			out = append(out, v)
		}
		return out
	}

	t.Run("yields top to bottom (LIFO)", func(t *testing.T) {
		s := New[int]()
		s.Push(1)
		s.Push(2)
		s.Push(3)

		if got, want := collect(s), []int{3, 2, 1}; !equalInts(got, want) {
			t.Errorf("All() = %v, want %v", got, want)
		}
	})

	t.Run("empty stack yields nothing", func(t *testing.T) {
		s := New[int]()
		if got := collect(s); len(got) != 0 {
			t.Errorf("All() over empty = %v, want no elements", got)
		}
	})

	t.Run("does not modify the stack", func(t *testing.T) {
		s := New[int]()
		s.Push(1)
		s.Push(2)

		_ = collect(s)

		if got := s.Len(); got != 2 {
			t.Errorf("Len() after All() = %d, want 2 (iteration must not consume)", got)
		}
		if v, ok := s.Pop(); !ok || v != 2 {
			t.Errorf("Pop() after All() = (%d, %v), want (2, true)", v, ok)
		}
	})

	t.Run("break stops iteration early and leaves stack intact", func(t *testing.T) {
		s := New[int]()
		s.Push(1)
		s.Push(2)
		s.Push(3)

		var got []int
		for v := range s.All() {
			got = append(got, v)
			if len(got) == 2 {
				break
			}
		}

		if want := []int{3, 2}; !equalInts(got, want) {
			t.Errorf("All() with early break = %v, want %v", got, want)
		}
		if l := s.Len(); l != 3 {
			t.Errorf("Len() after break = %d, want 3", l)
		}
	})
}

func TestStack_Clear(t *testing.T) {
	s := NewWithCap[int](8)
	for i := range 5 {
		s.Push(i)
	}

	s.Clear()

	if got := s.Len(); got != 0 {
		t.Errorf("Len() after Clear() = %d, want 0", got)
	}
	if !s.IsEmpty() {
		t.Errorf("IsEmpty() after Clear() = false, want true")
	}
	if got := s.Cap(); got != 0 {
		t.Errorf("Cap() after Clear() = %d, want 0 (backing array released)", got)
	}
	if v, ok := s.Pop(); ok {
		t.Errorf("Pop() after Clear() = (%d, true), want ok = false", v)
	}

	// stack stays usable after Clear
	s.Push(99)
	if v, ok := s.Pop(); !ok || v != 99 {
		t.Errorf("Pop() after reuse = (%d, %v), want (99, true)", v, ok)
	}
}

func TestStack_Reset(t *testing.T) {
	t.Run("empties stack but preserves capacity", func(t *testing.T) {
		s := New[int]()
		for i := range 10 {
			s.Push(i)
		}
		capBefore := s.Cap()

		s.Reset()

		if got := s.Len(); got != 0 {
			t.Errorf("Len() after Reset() = %d, want 0", got)
		}
		if !s.IsEmpty() {
			t.Errorf("IsEmpty() after Reset() = false, want true")
		}
		if got := s.Cap(); got != capBefore {
			t.Errorf("Cap() after Reset() = %d, want %d (capacity must be preserved)", got, capBefore)
		}
	})

	t.Run("zeroes backing array to release references", func(t *testing.T) {
		s := New[*int]()
		a, b, c := 1, 2, 3
		s.Push(&a)
		s.Push(&b)
		s.Push(&c)

		s.Reset()

		// white-box: the array is retained, so inspect every slot up to cap.
		full := s.data[:cap(s.data)]
		for i, p := range full {
			if p != nil {
				t.Errorf("backing slot %d = %p, want nil (Reset must zero elements)", i, p)
			}
		}
	})

	t.Run("stack stays usable after Reset", func(t *testing.T) {
		s := New[int]()
		s.Push(1)
		s.Reset()
		s.Push(42)

		if v, ok := s.Pop(); !ok || v != 42 {
			t.Errorf("Pop() after reuse = (%d, %v), want (42, true)", v, ok)
		}
		if v, ok := s.Pop(); ok {
			t.Errorf("Pop() on emptied stack = (%d, true), want ok = false", v)
		}
	})
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
