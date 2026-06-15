// Package queue provides a generic, slice-backed FIFO (first-in, first-out) queue.
//
// The zero value is not usable; create a queue with [New] or [NewWithCap].
// Elements are added to the back with [Queue.Push] and removed from the front
// with [Queue.Pop]; [Queue.Peek] inspects the front without removing it.
//
// Push and Pop run in amortized O(1) time (Push may occasionally reallocate
// the backing array as it grows). Peek, Len, Cap and IsEmpty run in O(1).
//
// A Queue is not safe for concurrent use: it performs no internal locking.
// Callers that share a queue across goroutines must provide their own
// synchronization.
package queue

// Queue is a generic FIFO queue backed by a slice. The front of the queue is the
// oldest element still present, i.e. the next one to be removed. A Queue must be
// created with [New] or [NewWithCap]; the zero value is not usable.
//
// Queue is not safe for concurrent use.
type Queue[T any] struct {
	data []T
}

// New returns an empty queue with no preallocated capacity.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// NewWithCap returns an empty queue with a backing array preallocated for at
// least n elements. Use it when the maximum size is known up front to avoid
// reallocations while pushing.
func NewWithCap[T any](n int) *Queue[T] {
	return &Queue[T]{data: make([]T, 0, n)}
}

// Len returns the number of elements currently in the queue.
func (q *Queue[T]) Len() int {
	return len(q.data)
}

// Cap returns the capacity of the queue's backing array measured from the
// current front, i.e. the number of elements it can hold before the next
// reallocation. Because [Queue.Pop] advances the front past the consumed slot,
// Cap shrinks as elements are popped.
func (q *Queue[T]) Cap() int {
	return cap(q.data)
}

// IsEmpty reports whether the queue has no elements.
func (q *Queue[T]) IsEmpty() bool {
	return len(q.data) == 0
}

// Peek returns the front element (the oldest, next to be removed) without
// removing it. The boolean result is false when the queue is empty, in which
// case the returned value is the zero value of T.
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	return q.data[0], true
}

// Push adds v to the back of the queue. It runs in amortized O(1) time; the
// backing array may be reallocated to grow.
func (q *Queue[T]) Push(v T) {
	q.data = append(q.data, v)
}

// Pop removes and returns the front element (the oldest, next to be removed).
// The boolean result is false when the queue is empty, in which case the
// returned value is the zero value of T. The vacated slot is zeroed so it no
// longer retains a reference to the popped element, then the front is advanced.
// It runs in O(1) time. Note that advancing the front reduces the capacity
// reported by [Queue.Cap].
func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if len(q.data) == 0 {
		return zero, false
	}
	v := q.data[0]
	q.data[0] = zero
	q.data = q.data[1:]
	return v, true
}

// Clear removes all elements and releases the backing array, so its memory
// becomes eligible for garbage collection. After Clear, Cap reports 0. To empty
// the queue while keeping its capacity for reuse, use [Queue.Reset] instead.
func (q *Queue[T]) Clear() {
	q.data = nil
}

// Reset removes all elements but keeps the backing array for reuse, so Cap is
// preserved. The elements are zeroed so the array no longer retains references
// to them. To also release the backing array, use [Queue.Clear] instead.
func (q *Queue[T]) Reset() {
	clear(q.data)
	q.data = q.data[:0]
}
