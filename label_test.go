package ctk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/kckrinke/go-cdk"
)

func TestLabel(t *testing.T) {
	Convey("Testing Labels", t, func() {
		Convey("basics: justification", func() {
			l := NewLabel("test")
			So(l, ShouldNotBeNil)
			So(l.GetTheme().String(), ShouldEqual, cdk.DefaultColorTheme.String())
			l.SetAllocation(cdk.MakeRectangle(10, 1))
			l.SetOrigin(0, 0)
			l.Resize()
			l.Show()
			canvas := cdk.NewCanvas(cdk.MakePoint2I(0, 0), cdk.MakeRectangle(10, 1), cdk.DefaultMonoTheme.Content.Normal)
			size := canvas.GetSize()
			So(size.W, ShouldEqual, 10)
			So(size.H, ShouldEqual, 1)
			So(l.Draw(canvas), ShouldEqual, cdk.EVENT_STOP)
			// canvas.ForEach(func(x, y int, cell cdk.TextCell) cdk.EventFlag {
			// 	So(y, ShouldEqual, 0)
			// 	switch x {
			// 	case 0:
			// 		So(cell.Value(), ShouldEqual, 't')
			// 	case 1:
			// 		So(cell.Value(), ShouldEqual, 'e')
			// 	case 2:
			// 		So(cell.Value(), ShouldEqual, 's')
			// 	case 3:
			// 		So(cell.Value(), ShouldEqual, 't')
			// 	default:
			// 		So(cell.IsSpace(), ShouldEqual, true)
			// 	}
			// 	return cdk.EVENT_PASS
			// })
			l.SetJustify(cdk.JUSTIFY_RIGHT)
			canvas = cdk.NewCanvas(cdk.MakePoint2I(0, 0), cdk.MakeRectangle(10, 1), cdk.DefaultMonoTheme.Content.Normal)
			size = canvas.GetSize()
			So(size.W, ShouldEqual, 10)
			So(size.H, ShouldEqual, 1)
			So(l.Draw(canvas), ShouldEqual, cdk.EVENT_STOP)
			// canvas.ForEach(func(x, y int, cell cdk.TextCell) cdk.EventFlag {
			// 	So(y, ShouldEqual, 0)
			// 	switch x {
			// 	case 6:
			// 		So(cell.Value(), ShouldEqual, 't')
			// 	case 7:
			// 		So(cell.Value(), ShouldEqual, 'e')
			// 	case 8:
			// 		So(cell.Value(), ShouldEqual, 's')
			// 	case 9:
			// 		So(cell.Value(), ShouldEqual, 't')
			// 	default:
			// 		So(cell.IsSpace(), ShouldEqual, true)
			// 	}
			// 	return cdk.EVENT_PASS
			// })
			l.SetJustify(cdk.JUSTIFY_CENTER)
			canvas = cdk.NewCanvas(cdk.MakePoint2I(0, 0), cdk.MakeRectangle(10, 1), cdk.DefaultMonoTheme.Content.Normal)
			size = canvas.GetSize()
			So(size.W, ShouldEqual, 10)
			So(size.H, ShouldEqual, 1)
			So(l.Draw(canvas), ShouldEqual, cdk.EVENT_STOP)
			// canvas.ForEach(func(x, y int, cell cdk.TextCell) cdk.EventFlag {
			// 	So(y, ShouldEqual, 0)
			// 	switch x {
			// 	case 3:
			// 		So(cell.Value(), ShouldEqual, 't')
			// 	case 4:
			// 		So(cell.Value(), ShouldEqual, 'e')
			// 	case 5:
			// 		So(cell.Value(), ShouldEqual, 's')
			// 	case 6:
			// 		So(cell.Value(), ShouldEqual, 't')
			// 	default:
			// 		So(cell.IsSpace(), ShouldEqual, true)
			// 	}
			// 	return cdk.EVENT_PASS
			// })
		})
	})
}
