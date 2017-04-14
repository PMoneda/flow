package gonnie

import "testing"

func TestShouldCreateAWorkFlow(t *testing.T) {
	c := NewContext()
	r := NewRoute(c)
	r = r.From("direct://hello").To("http://google.com").To("file://destino?name=ggg.txt")

}
