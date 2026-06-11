# argos

Generic, idiomatic data structures for Go — small, allocation-friendly, and
built on the Go 1.23+ standard library, including range-over-func iterators.

`argos` is a growing collection of generic container types. Each package is
self-contained, dependency-free, and documented with runnable examples.

## Requirements

Go 1.23+ (for range-over-func iterators). The module targets Go 1.26.

## Install

```go
import "argos/stack"
```

## Packages

| Package            | Description                                                   | Docs                          |
| ------------------ | ------------------------------------------------------------- | ----------------------------- |
| [`stack`](./stack) | Generic, slice-backed LIFO stack with non-consuming iteration. | [README](./stack/README.md) |

More containers are on the way.

## Quick start

```go
package main

import (
	"fmt"

	"argos/stack"
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

## Design goals

- **Generic** — type-safe containers via Go generics, no `interface{}` boxing.
- **Idiomatic** — small APIs that follow standard-library conventions, such as
  the `(value, ok)` contract and [`iter.Seq`](https://pkg.go.dev/iter#Seq)
  iterators.
- **Allocation-friendly** — capacity preallocation where it matters, with
  explicit control over when backing memory is released.
- **Dependency-free** — standard library only.

## Documentation

Each package ships a detailed README and verified `example_test.go` files that
render on godoc:

```sh
go doc argos/stack          # package overview
go doc argos/stack Stack    # a type and its methods
```

## License

[MIT](./LICENSE) © Nergous
