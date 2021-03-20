package ctk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/kckrinke/go-cdk"
)

func TestVScrollbar(t *testing.T) {
	Convey("Testing raw scrollbar", t, func() {
		s := &CScrollbar{}
		Convey("default to vertical", func() {
			So(s.orientation, ShouldEqual, 0)
			So(s.Init(), ShouldEqual, false)
			So(s.Init(), ShouldEqual, true)
			So(s.orientation, ShouldEqual, cdk.ORIENTATION_VERTICAL)
		})
	})
	Convey("Testing vertical scrollbars", t, func() {
		Convey("basic checks", func() {
			vs := NewVScrollbar()
			So(vs, ShouldNotBeNil)
			So(vs.GetHasBackwardStepper(), ShouldEqual, true)
			vs.SetHasBackwardStepper(false)
			So(vs.GetHasBackwardStepper(), ShouldEqual, false)
		})
	})
}
