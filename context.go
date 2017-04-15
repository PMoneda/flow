package gonnie

import (
	"fmt"
)

// StackExchange stack for messages
type StackExchange []*Exchange

//StringStack is a stack based on []string
type StringStack []string

// Push a message to stack
func (s *StringStack) Push(ext string) {
	*s = append(*s, ext)
}

// Pop a message from stack
func (s *StringStack) Pop() string {
	if len(*s) == 0 {
		return ""
	}
	var x string
	x, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return x
}

// Print log from stack
func (s StringStack) Print() {
	for _, str := range s {
		fmt.Println(str)
	}
}

// Clear stack
func (s StringStack) Clear() {
	str := s.Pop()
	for str != "" {
		str = s.Pop()
	}
}

// Push a message to stack
func (s *StackExchange) Push(ext *Exchange) {
	*s = append(*s, ext)
}

// Pop a message from stack
func (s *StackExchange) Pop() *Exchange {
	if len(*s) == 0 {
		return nil
	}
	var x *Exchange
	x, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return x
}

// Top return a message from top of stack
func (s StackExchange) Top() *Exchange {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}

// Context keep message in the flow
type Context struct {
	stack StackExchange
	log   StringStack
}

// PushMessage push message to Context
func (c *Context) PushMessage() {
	new := NewExchange()
	top := c.stack.Top()
	if top != nil {
		new.in = top.out
		new.inHead = top.outHead
	}
	c.stack.Push(new)
}

// PopMessage pop message from Context
func (c *Context) PopMessage() *Exchange {
	return c.stack.Pop()
}

// GetMessage from context
func (c *Context) GetMessage() *Exchange {
	if len(c.stack) == 0 {
		c.PushMessage()
	}
	return c.stack.Top()
}

// GetLog from context
func (c *Context) GetLog() StringStack {
	return c.log
}

// NewContext returns a empty context
func NewContext() *Context {
	c := Context{stack: make(StackExchange, 0, 0), log: make(StringStack, 0, 0)}
	c.stack = append(c.stack, NewExchange())
	return &c
}
