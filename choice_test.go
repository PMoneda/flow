package flow

import (
	"fmt"
	"testing"
)

func TestChoice(t *testing.T) {
	p := NewPipe().From("direct:a")
	p = p.SetBody("Hello").SetHeader("status", "300")
	ch := p.Choice()
	ch = ch.When(Header("status").IsEqualTo("200")).To("print://?msg=200")
	ch = ch.When(Header("status").IsEqualTo("400")).To("print://?msg=400")
	b := ch.Otherwise().Body()
	fmt.Println(b)
}
