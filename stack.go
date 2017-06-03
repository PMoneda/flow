package flow

// Stack for message channel
type Stack []interface{}

// Push a message to stack
func (s *Stack) Push(ext interface{}) {
	*s = append(*s, ext)
}

//Clear stack
func (s *Stack) Clear() {
	s = nil
}

// Pop a message from stack
func (s *Stack) Pop() interface{} {
	p := *s
	if len(p) == 0 {
		return nil
	}
	var x interface{}
	x, p = p[len(p)-1], p[:len(p)-1]
	return x
}

// Top return a message from top of stack
func (s Stack) Top() interface{} {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}
