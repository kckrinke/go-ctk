package ctk

import (
	"github.com/kckrinke/go-cdk"
)

const (
	TypeVBox cdk.CTypeTag = "ctk-v-box"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeVBox, func() interface{} { return MakeVBox() })
}

// Basic vbox interface
type VBox interface {
	Box

	Init() bool
}

type CVBox struct {
	CBox
}

func MakeVBox() *CVBox {
	return NewVBox(false, 0)
}

func NewVBox(homogeneous bool, spacing int) *CVBox {
	b := new(CVBox)
	b.Init()
	b.SetHomogeneous(homogeneous)
	b.SetSpacing(spacing)
	return b
}

func (b *CVBox) Init() bool {
	if b.InitTypeItem(TypeVBox, b) {
		return true
	}
	b.CBox.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(PARENT_SENSITIVE)
	b.SetFlags(APP_PAINTABLE)
	b.SetOrientation(cdk.ORIENTATION_VERTICAL)
	return false
}
