package util

const arraySize = 4096

type Stack struct {
	top  int
	data [arraySize]interface{}
}

func (s *Stack) Push(i interface{}) bool {
	if s.top == len(s.data) {
		return false
	}
	s.data[s.top] = i
	s.top += 1
	return true
}

func (s *Stack) Pop() (interface{}, bool) {
	if s.top == 0 {
		return 0, false
	}
	i := s.data[s.top-1]
	s.top -= 1
	return i, true
}

func (s *Stack) Peek() interface{} {
	return s.data[s.top-1]
}

func (s *Stack) Get() []interface{} {
	return s.data[:s.top]
}

func (s *Stack) IsEmpty() bool {
	return s.top == 0
}

func (s *Stack) Empty() {
	s.top = 0
}
