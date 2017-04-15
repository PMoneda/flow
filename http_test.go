package gonnie

import "testing"

func TestShouldGetRequestFromUrl(t *testing.T) {
	ctx := NewContext()
	route := NewRoute(ctx)
	route.From("https://?url=http://correiosapi.apphb.com/cep/76873274&method=GET")
	route.Processor(func(ext *Exchange) {
		ext.GetOut().WriteString("CHANGE" + ext.GetIn().String())
	})
}
