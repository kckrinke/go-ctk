package ctk

import (
	"github.com/kckrinke/go-cdk"
)

const TypeHButtonBox cdk.CTypeTag = "ctk-h-button-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeHButtonBox, func() interface{} { return MakeHButtonBox() })
}

type HButtonBox interface {
	ButtonBox

	Init() bool
}

type CHButtonBox struct {
	CButtonBox
}

func MakeHButtonBox() *CHButtonBox {
	return NewHButtonBox(false, 0)
}

func NewHButtonBox(homogeneous bool, spacing int) *CHButtonBox {
	b := new(CHButtonBox)
	b.Init()
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

func (b *CHButtonBox) Init() bool {
	if b.InitTypeItem(TypeHButtonBox, b) {
		return true
	}
	b.CButtonBox.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE)
	b.SetFlags(APP_PAINTABLE)
	b.SetOrientation(cdk.ORIENTATION_HORIZONTAL)
	return false
}
