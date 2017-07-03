package flow

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldSetBasicHeaderAndBody(t *testing.T) {
	Convey("Should set header to message", t, func() {
		f := NewFlow().From("direct://")
		f.SetHeader("a", "b")
		So(f.GetHeader().Get("a"), ShouldEqual, "b")
	})

	Convey("Should set body to message", t, func() {
		f := NewFlow().From("direct://", 1)
		f.SetBody(2)
		So(f.GetBody(), ShouldEqual, 2)
	})
}
