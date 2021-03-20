package ctk

import (
	"github.com/kckrinke/go-cdk"
)

const (
	TypeHScrollbar cdk.CTypeTag = "ctk-h-scrollbar"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeHScrollbar, func() interface{} { return MakeHScrollbar() })
}

type HScrollbar interface {
	Scrollbar

	Init() (already bool)
}

type CHScrollbar struct {
	CScrollbar
}

func MakeHScrollbar() *CHScrollbar {
	return NewHScrollbar()
}

func NewHScrollbar() *CHScrollbar {
	s := &CHScrollbar{}
	s.orientation = cdk.ORIENTATION_HORIZONTAL
	s.Init()
	return s
}

func (s *CHScrollbar) Init() (already bool) {
	if s.InitTypeItem(TypeHScrollbar, s) {
		return true
	}
	s.CScrollbar.Init()
	s.SetFlags(SENSITIVE | PARENT_SENSITIVE)
	s.SetFlags(CAN_FOCUS)
	s.SetFlags(APP_PAINTABLE)
	s.SetTheme(DefaultColorScrollbarTheme)
	return false
}
