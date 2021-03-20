package ctk

import (
	"github.com/kckrinke/go-cdk"
)

const TypeFakeWindow cdk.CTypeTag = "ctk-fake-window"

type WithFakeWindowFn = func(w Window)

type CFakeWindow struct {
	CWindow
}

func (f *CFakeWindow) Init() (already bool) {
	if f.InitTypeItem(TypeFakeWindow, f) {
		return true
	}
	f.SetAllocation(cdk.MakeRectangle(80, 24))
	f.SetTheme(cdk.DefaultNilTheme)
	return false
}

func (f *CFakeWindow) FakeDraw() (canvas cdk.Canvas, flag cdk.EventFlag) {
	canvas = cdk.NewCanvas(cdk.Point2I{}, f.GetAllocation(), f.GetTheme().Content.Normal)
	flag = f.Draw(canvas)
	return
}

func WithFakeWindow(fn WithFakeWindowFn) func() {
	fakeWindow := new(CFakeWindow)
	fakeWindow.Init()
	fakeWindow.SetTheme(cdk.DefaultMonoTheme)
	return func() {
		fn(fakeWindow)
	}
}

func WithFakeWindowOptions(w, h int, theme cdk.Theme, fn WithFakeWindowFn) func() {
	fakeWindow := new(CFakeWindow)
	fakeWindow.Init()
	fakeWindow.SetAllocation(cdk.MakeRectangle(w, h))
	fakeWindow.SetTheme(theme)
	return func() {
		fn(fakeWindow)
	}
}
