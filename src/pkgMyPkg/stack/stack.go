package stack

import "errors"

type stack struct {
	top int
	arr []interface{}
}

func (s *stack) Len() int {
	return len(s.arr)
}

func (s *stack) IsEmpty() bool {
	return s.top == 0
}

func (s *stack) Push(x interface{}) error {
	if s.top == cap(s.arr) {
		return errors.New("full stack")
	}

	s.arr[s.top] = x
	s.top++
	return nil
}

func (s *stack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("empty stack")
	}
	s.top--
	return s.arr[s.top], nil
}

func (s *stack) Top() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("empty stack")
	}
	return s.arr[s.top-1], nil
}

func New(maxNum int) *stack {
	s := new(stack)
	s.top = 0
	s.arr = make([]interface{}, maxNum)
	return s
}
