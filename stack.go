package hocon

import (
	"errors"
)

type stack struct {
	head    int
	content []valueSetter
}

func NewStack() *stack {
	return &stack{
		head:    -1,
		content: make([]valueSetter, 1),
	}
}

func (s *stack) Push(p valueSetter) {
	s.head++
	length := len(s.content)
	if s.head >= length {
		tmp := make([]valueSetter, length*2)
		copy(tmp, s.content)
		for i := 0; i < length; i++ {
			s.content[i] = nil
		}
		s.content = tmp
	}
	s.content[s.head] = p
}

func (s *stack) Pop() (v valueSetter, err error) {
	if s.head < 0 {
		v = nil
		err = errors.New("Pop on the empty stack")
	} else {
		v = s.content[s.head]
		s.content[s.head] = nil
		s.head--
	}
	return v, err
}

func (s *stack) Peek() (v valueSetter, err error) {
	if s.head < 0 {
		v = nil
		err = errors.New("Peek on the empty stack")
	} else {
		v = s.content[s.head]
	}
	return v, err
}
