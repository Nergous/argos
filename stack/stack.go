// Package stack provides a generic, slice-backed LIFO (last-in, first-out) stack.
//
// The zero value is not usable; create a stack with [New] or [NewWithCap].
// Elements are added with [Stack.Push] and removed from the top with
// [Stack.Pop]; [Stack.Peek] inspects the top without removing it.
//
// Push and Pop run in amortized O(1) time (Push may occasionally reallocate
// the backing array as it grows). Peek, Len, Cap and IsEmpty run in O(1).
//
// A Stack is not safe for concurrent use: it performs no internal locking.
// Callers that share a stack across goroutines must provide their own
// synchronization.
package stack

import "iter"

// Stack is a generic LIFO stack backed by a slice. The top of the stack is the
// most recently pushed element. A Stack must be created with [New] or
// [NewWithCap]; the zero value is not usable.
//
// Stack is not safe for concurrent use.
type Stack[T any] struct {
	data []T
}

// New returns an empty stack with no preallocated capacity.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// NewWithCap returns an empty stack with a backing array preallocated for at
// least n elements. Use it when the maximum size is known up front to avoid
// reallocations while pushing.
func NewWithCap[T any](n int) *Stack[T] {
	return &Stack[T]{data: make([]T, 0, n)}
}

// Len returns the number of elements currently on the stack.
func (s *Stack[T]) Len() int {
	return len(s.data)
}

// Cap returns the capacity of the stack's backing array, i.e. the number of
// elements it can hold before the next reallocation.
func (s *Stack[T]) Cap() int {
	return cap(s.data)
}

// IsEmpty reports whether the stack has no elements.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) < 1
}

// Peek returns the top element without removing it. The boolean result is false
// when the stack is empty, in which case the returned value is the zero value
// of T.
func (s *Stack[T]) Peek() (T, bool) {
	l := len(s.data)
	if l < 1 {
		var zero T
		return zero, false
	}

	return s.data[l-1], true
}

// All returns an iterator over the stack's elements from top to bottom (LIFO
// order), without consuming them; the stack is left unchanged. The returned
// [iter.Seq] is a Go 1.23 range-over-func iterator, so it can be used directly
// in a range loop. Breaking out of the loop early stops iteration cleanly.
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(s.data) - 1; i >= 0; i-- {
			if !yield(s.data[i]) {
				return
			}
		}
	}
}

// Push adds v to the top of the stack. It runs in amortized O(1) time; the
// backing array may be reallocated to grow.
func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

// Pop removes and returns the top element. The boolean result is false when the
// stack is empty, in which case the returned value is the zero value of T. The
// vacated slot is zeroed so it no longer retains a reference to the popped
// element. It runs in O(1) time.
func (s *Stack[T]) Pop() (T, bool) {
	l := len(s.data) - 1
	var zero T

	if l < 0 {
		return zero, false
	}

	value := s.data[l]
	s.data[l] = zero
	s.data = s.data[:l]

	return value, true
}

// Clear removes all elements and releases the backing array, so its memory
// becomes eligible for garbage collection. After Clear, Cap reports 0. To empty
// the stack while keeping its capacity for reuse, use [Stack.Reset] instead.
func (s *Stack[T]) Clear() {
	s.data = nil
}

// Reset removes all elements but keeps the backing array for reuse, so Cap is
// preserved. The elements are zeroed so the array no longer retains references
// to them. To also release the backing array, use [Stack.Clear] instead.
func (s *Stack[T]) Reset() {
	clear(s.data)
	s.data = s.data[:0]
}
