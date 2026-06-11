# stack

A generic, slice-backed **LIFO** (last-in, first-out) stack for Go.

It is part of the `argos` data-structures library. The stack is small,
allocation-friendly, and uses Go 1.23 range-over-func iterators
([`iter.Seq`](https://pkg.go.dev/iter#Seq)) for non-consuming iteration.

## Install

```go
import "github.com/Nergous/argos/stack"
```

Requires Go 1.23+ (for range-over-func). The module targets Go 1.26.

## Usage

```go
package main

import (
	"fmt"

	"github.com/Nergous/argos/stack"
)

func main() {
	s := stack.New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	for !s.IsEmpty() {
		v, _ := s.Pop()
		fmt.Println(v) // 3, 2, 1
	}
}
```

Preallocate when the maximum size is known up front to avoid reallocations
while pushing:

```go
s := stack.NewWithCap[string](3)
s.Push("a")
s.Push("b")
s.Push("c")
fmt.Println(s.Len(), s.Cap()) // 3 3
```

## API reference

| Method / Function       | Signature                              | Description                                                                                         | Complexity      |
| ----------------------- | -------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------- |
| `New`                   | `func New[T any]() *Stack[T]`           | Creates an empty stack with no preallocated capacity.                                                | O(1)            |
| `NewWithCap`            | `func NewWithCap[T any](n int) *Stack[T]` | Creates an empty stack with capacity preallocated for at least `n` elements.                        | O(n)            |
| `Len`                   | `func (s *Stack[T]) Len() int`          | Number of elements currently on the stack.                                                          | O(1)            |
| `Cap`                   | `func (s *Stack[T]) Cap() int`          | Capacity of the backing array before the next reallocation.                                         | O(1)            |
| `IsEmpty`               | `func (s *Stack[T]) IsEmpty() bool`     | Reports whether the stack has no elements.                                                           | O(1)            |
| `Peek`                  | `func (s *Stack[T]) Peek() (T, bool)`   | Returns the top element without removing it; `ok` is `false` and the value is the zero value if empty. | O(1)            |
| `All`                   | `func (s *Stack[T]) All() iter.Seq[T]`  | Iterator over elements from top to bottom (LIFO) without consuming them.                             | O(1) per step   |
| `Push`                  | `func (s *Stack[T]) Push(v T)`          | Adds `v` to the top of the stack.                                                                   | amortized O(1)  |
| `Pop`                   | `func (s *Stack[T]) Pop() (T, bool)`    | Removes and returns the top element; `ok` is `false` and the value is the zero value if empty.       | O(1)            |
| `Clear`                 | `func (s *Stack[T]) Clear()`            | Removes all elements **and releases** the backing array (`Cap` drops to 0).                         | O(1)            |
| `Reset`                 | `func (s *Stack[T]) Reset()`            | Removes all elements but **keeps** the backing array for reuse (`Cap` preserved).                   | O(n)            |

### The `(T, bool)` contract

`Peek` and `Pop` return two values. The boolean reports whether the operation
found an element: it is `true` on success and `false` when the stack is empty.
When it is `false`, the first return is the zero value of `T`.

```go
v, ok := s.Pop()
if !ok {
	// stack was empty; v is the zero value
}
```

### Non-consuming iteration with `All`

`All` returns an [`iter.Seq[T]`](https://pkg.go.dev/iter#Seq) — a Go 1.23
range-over-func iterator — that yields elements from **top to bottom** (LIFO
order). It does not modify the stack, and breaking out of the loop early stops
iteration cleanly:

```go
for v := range s.All() {
	fmt.Println(v)
	if v == target {
		break // safe; the stack is untouched
	}
}
fmt.Println(s.Len()) // unchanged
```

## `Clear` vs `Reset`

Both empty the stack and zero the removed elements (so the backing array no
longer retains references to them), but they differ in what happens to the
allocated memory:

| Operation | Empties stack | Backing array        | `Cap` afterwards | Use when                                                  |
| --------- | ------------- | -------------------- | ---------------- | --------------------------------------------------------- |
| `Clear`   | yes           | released (eligible for GC) | `0`        | You are done with the stack, or want to free its memory.  |
| `Reset`   | yes           | retained for reuse   | unchanged        | You will refill the stack and want to avoid reallocating. |

```go
s := stack.NewWithCap[int](8)
s.Push(1)
s.Push(2)

s.Reset()
fmt.Println(s.Len(), s.Cap()) // 0 8  — capacity kept for reuse

s.Clear()
fmt.Println(s.Len(), s.Cap()) // 0 0  — backing array released
```

## Concurrency

A `Stack` is **not safe for concurrent use**: it performs no internal locking.
If a stack is shared across goroutines, the caller must provide its own
synchronization (for example, a `sync.Mutex`).

## More examples and docs

Runnable, verified examples live in
[`example_test.go`](./example_test.go) — including bracket-balancing — and are
rendered alongside the API on the godoc page.

View the documentation locally:

```sh
go doc github.com/Nergous/argos/stack          # package overview
go doc github.com/Nergous/argos/stack Stack    # the Stack type and its methods
```

The same comments render on pkg.go.dev-style godoc.
