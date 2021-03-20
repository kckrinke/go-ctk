package ctk

import (
	"github.com/kckrinke/go-cdk"
)

type Drawable interface {
	Hide()
	Show()
	ShowAll()
	IsVisible() bool
	HasPoint(p *cdk.Point2I) bool
	GetWidgetAt(p *cdk.Point2I) (instance interface{})
	GetSizeRequest() (size cdk.Rectangle)
	SetSizeRequest(x, y int)
	GetTheme() (theme cdk.Theme)
	SetTheme(theme cdk.Theme)
	GetThemeRequest() (theme cdk.Theme)
	GetOrigin() (origin cdk.Point2I)
	SetOrigin(x, y int)
	GetAllocation() (alloc cdk.Rectangle)
	SetAllocation(alloc cdk.Rectangle)
	Invalidate() cdk.EventFlag
	Resize() cdk.EventFlag
	Draw(canvas cdk.Canvas) cdk.EventFlag
}
