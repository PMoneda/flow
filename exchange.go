package gonnie

import "bytes"

//Exchange is a middleware message used between process
type Exchange struct {
	in  *bytes.Buffer
	out *bytes.Buffer
}

// GetIn returns input message
func (e *Exchange) GetIn() *bytes.Buffer {
	return e.in
}

// GetOut returns output buffer
func (e *Exchange) GetOut() *bytes.Buffer {
	return e.out
}

// NewExchange creates new exchange message
func NewExchange() *Exchange {
	e := Exchange{
		in:  bytes.NewBuffer(nil),
		out: bytes.NewBuffer(nil),
	}
	return &e
}
