package main

import (
	"errors"
)

type stakElement = int

var emptyElement stakElement

type stack struct {
	top int // 指向下一个空位置
	arr []stakElement
}

func (s *stack) Len() int {
	return len(s.arr)
}

func (s *stack) IsEmpty() bool {
	return s.top == 0
}

func (s *stack) Push(x stakElement) error {
	s.arr = append(s.arr, x)
	s.top++
	return nil
}

func (s *stack) Pop() (stakElement, error) {
	if s.IsEmpty() {
		return emptyElement, errors.New("empty stack")
	}
	e := s.arr[s.top-1]
	s.top--
	s.arr = s.arr[:s.top]

	return e, nil
}

func (s *stack) Top() (stakElement, error) {
	if s.IsEmpty() {
		return emptyElement, errors.New("empty stack")
	}
	return s.arr[s.top-1], nil
}

func NewStack(capNum int) *stack {
	s := new(stack)
	s.top = 0
	s.arr = make([]stakElement, 0, capNum)
	return s
}
