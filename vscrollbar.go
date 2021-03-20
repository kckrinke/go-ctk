package ctk

import (
	"github.com/kckrinke/go-cdk"
)

const (
	TypeVScrollbar cdk.CTypeTag = "ctk-v-scrollbar"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeVScrollbar, func() interface{} { return MakeVScrollbar() })
}

type VScrollbar interface {
	Scrollbar

	Init() (already bool)
}

type CVScrollbar struct {
	CScrollbar
}

func MakeVScrollbar() *CVScrollbar {
	return NewVScrollbar()
}

func NewVScrollbar() *CVScrollbar {
	v := &CVScrollbar{}
	v.orientation = cdk.ORIENTATION_VERTICAL
	v.Init()
	return v
}

func (v *CVScrollbar) Init() (already bool) {
	if v.InitTypeItem(TypeVScrollbar, v) {
		return true
	}
	v.CScrollbar.Init()
	v.SetFlags(SENSITIVE | PARENT_SENSITIVE)
	v.SetFlags(CAN_FOCUS)
	v.SetFlags(APP_PAINTABLE)
	v.SetTheme(DefaultColorScrollbarTheme)
	return false
}
