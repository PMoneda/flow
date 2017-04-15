package gonnie

import "testing"

func TestShouldParseURL(t *testing.T) {
	u := "file://desktop.documents?name=teste.txt"
	ur, err := processURI(u)
	if err != nil {
		t.Fail()
	}

	if ur.protocol != "file" {
		t.Fail()
	}
}

func TestBaseFlow(t *testing.T) {
	ctx := NewContext()
	route := NewRoute(ctx)
	route.From("direct://route").Log("A").Log("B").Log("C")
	ctx.GetLog().Clear()
}
