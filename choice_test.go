package flow

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChoice(t *testing.T) {
	executeFlow := func(s string) string {
		p := NewPipe().From("direct:a")
		p = p.SetBody("Hello").SetHeader("status", s)
		ch := p.Choice()
		ch = ch.When(Header("status").IsEqualTo("200")).To("set://?prop=body", "200")
		ch = ch.When(Header("status").IsEqualTo("400")).To("set://?prop=body", "400")
		ch = ch.Otherwise().To("set://?prop=body", "500")
		b := p.GetBody()
		switch t := b.(type) {
		case string:
			return t
		default:
			return "-1"

		}
	}
	Convey("Should decide test coditional flows", t, func() {

		So(executeFlow("200"), ShouldEqual, "200")
		So(executeFlow("400"), ShouldEqual, "400")
		So(executeFlow("300"), ShouldEqual, "500")

	})

}
