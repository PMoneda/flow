package gonnie

// StackExchange stack for messages
type StackExchange []*Exchange

// Push a message to stack
func (s StackExchange) Push(ext *Exchange) {
	s = append(s, ext)
}

// Pop a message from stack
func (s StackExchange) Pop() *Exchange {
	var x *Exchange
	x, s = s[len(s)-1], s[:len(s)-1]
	return x
}

// Top return a message from top of stack
func (s StackExchange) Top() *Exchange {
	return s[len(s)-1]
}

// Context keep message in the flow
type Context struct {
	msg StackExchange
}

// PushMessage push message to Context
func (c *Context) PushMessage(msg *Exchange) {
	c.msg.Push(msg)
}

// GetMessage from context
func (c *Context) GetMessage() *Exchange {
	return c.msg.Top()
}

// NewContext returns a empty context
func NewContext() *Context {
	c := Context{msg: make(StackExchange, 0, 1)}
	c.msg[0] = NewExchange()
	return &c
}
