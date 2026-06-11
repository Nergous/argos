package stack_test

import (
	"fmt"

	"argos/stack"
)

// A stack is LIFO (last in, first out): elements come out in the reverse
// order they went in.
func Example() {
	s := stack.New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	for !s.IsEmpty() {
		v, _ := s.Pop()
		fmt.Println(v)
	}
	// Output:
	// 3
	// 2
	// 1
}

// NewWithCap preallocates the backing array. Use it when the final size is
// known up front to avoid reallocations during Push.
func ExampleNewWithCap() {
	s := stack.NewWithCap[string](3)
	s.Push("a")
	s.Push("b")
	s.Push("c")

	fmt.Println(s.Len(), s.Cap())
	// Output: 3 3
}

// Pop removes and returns the top element. The second value reports whether
// the stack was non-empty; on an empty stack it is false and the value is the
// zero value.
func ExampleStack_Pop() {
	s := stack.New[string]()
	s.Push("first")
	s.Push("last")

	v, ok := s.Pop()
	fmt.Println(v, ok)

	s.Pop() // removes "first"

	_, ok = s.Pop() // stack is now empty
	fmt.Println(ok)
	// Output:
	// last true
	// false
}

// Peek returns the top element without removing it.
func ExampleStack_Peek() {
	s := stack.New[int]()
	s.Push(10)
	s.Push(20)

	top, ok := s.Peek()
	fmt.Println(top, ok)
	fmt.Println("len unchanged:", s.Len())
	// Output:
	// 20 true
	// len unchanged: 2
}

// All returns a range-over-func iterator over the elements from top to bottom
// without consuming them. The stack is left intact after iteration.
func ExampleStack_All() {
	s := stack.New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	for v := range s.All() {
		fmt.Println(v)
	}
	fmt.Println("remaining:", s.Len())
	// Output:
	// 3
	// 2
	// 1
	// remaining: 3
}

// Reset empties the stack but keeps the allocated backing array, so it can be
// reused without a new allocation.
func ExampleStack_Reset() {
	s := stack.NewWithCap[int](8)
	s.Push(1)
	s.Push(2)

	s.Reset()
	fmt.Println(s.Len(), s.Cap())
	// Output: 0 8
}

// Clear, unlike Reset, releases the backing array: capacity drops back to zero
// and the memory becomes eligible for garbage collection.
func ExampleStack_Clear() {
	s := stack.NewWithCap[int](8)
	s.Push(1)
	s.Push(2)

	s.Clear()
	fmt.Println(s.Len(), s.Cap())
	// Output: 0 0
}

// A practical example: checking whether brackets are balanced. Push each
// opening bracket; on a closing one, pop and verify the pair matches.
func Example_balancedBrackets() {
	balanced := func(input string) bool {
		pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
		st := stack.New[rune]()

		for _, r := range input {
			switch r {
			case '(', '[', '{':
				st.Push(r)
			case ')', ']', '}':
				top, ok := st.Pop()
				if !ok || top != pairs[r] {
					return false
				}
			}
		}
		return st.IsEmpty()
	}

	fmt.Println(balanced("([]{})"))
	fmt.Println(balanced("([)]"))
	// Output:
	// true
	// false
}