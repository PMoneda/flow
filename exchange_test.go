package gonnie

import "testing"

func TestShouldCreateExchange(t *testing.T) {
	e := NewExchange()
	e.GetOut().WriteString("Hello")
	if e.GetOut().String() != "Hello" {
		t.Fail()
	}
}
