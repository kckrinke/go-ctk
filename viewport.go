package ctk

import (
	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
)

// CDK type-tag for Viewport objects
const TypeViewport cdk.CTypeTag = "ctk-viewport"

func init() {
	_ = cdk.TypesManager.AddType(TypeViewport, func() interface{} { return MakeViewport() })
}

// Viewport Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Viewport
// The Viewport widget acts as an adaptor class, implementing scrollability
// for child widgets that lack their own scrolling capabilities. Use Viewport
// to scroll child widgets such as Table, Box, and so on. If a widget has
// native scrolling abilities, such as TextView, TreeView or Iconview, it can
// be added to a ScrolledWindow with ContainerAdd. If a widget does
// not, you must first add the widget to a Viewport, then add the viewport to
// the scrolled window. The convenience function
// ScrolledWindowAddWithViewport does exactly this, so you can
// ignore the presence of the viewport.
type Viewport interface {
	Bin

	Init() (already bool)
	GetHAdjustment() *CAdjustment
	SetHAdjustment(adjustment *CAdjustment)
	GetVAdjustment() *CAdjustment
	SetVAdjustment(adjustment *CAdjustment)
	SetShadowType(shadowType ShadowType)
	GetShadowType() (value ShadowType)
	GetBinWindow() (value Window)
	GetViewWindow() (value Window)
	Invalidate() cdk.EventFlag
	Resize() cdk.EventFlag
	Draw(canvas cdk.Canvas) cdk.EventFlag
}

// The CViewport structure implements the Viewport interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Viewport objects
type CViewport struct {
	CBin

	canvas *cdk.CCanvas
}

// Default constructor for Viewport objects
func MakeViewport() *CViewport {
	return NewViewport(nil, nil)
}

func NewViewport(hAdjustment, vAdjustment *CAdjustment) *CViewport {
	v := new(CViewport)
	v.Init()
	v.SetHAdjustment(hAdjustment)
	v.SetVAdjustment(vAdjustment)
	return v
}

// Viewport object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Viewport instance
func (v *CViewport) Init() (already bool) {
	if v.InitTypeItem(TypeViewport, v) {
		return true
	}
	v.CBin.Init()
	v.flags = NULL_WIDGET_FLAG
	v.SetFlags(SENSITIVE | PARENT_SENSITIVE)
	v.SetFlags(APP_PAINTABLE)
	_ = v.InstallProperty(PropertyHAdjustment, cdk.StructProperty, true, nil)
	_ = v.InstallProperty(PropertyViewportShadowType, cdk.StructProperty, true, nil)
	_ = v.InstallProperty(PropertyVAdjustment, cdk.StructProperty, true, nil)
	v.canvas = nil
	return false
}

// Returns the horizontal adjustment of the viewport.
// Returns:
// 	the horizontal adjustment of viewport .
// 	[transfer none]
func (v *CViewport) GetHAdjustment() (adjustment *CAdjustment) {
	var ok bool
	if value, err := v.GetStructProperty(PropertyHAdjustment); err != nil {
		v.LogErr(err)
	} else if adjustment, ok = value.(*CAdjustment); !ok {
		v.LogError("value stored in %v property is not of *CAdjustment type: %v (%T)", PropertyHAdjustment, value, value)
	}
	return
}

// Returns the vertical adjustment of the viewport.
// Returns:
// 	the vertical adjustment of viewport .
// 	[transfer none]
func (v *CViewport) GetVAdjustment() (adjustment *CAdjustment) {
	var ok bool
	if value, err := v.GetStructProperty(PropertyVAdjustment); err != nil {
		v.LogErr(err)
	} else if adjustment, ok = value.(*CAdjustment); !ok {
		v.LogError("value stored in %v property is not of *CAdjustment type: %v (%T)", PropertyVAdjustment, value, value)
	}
	return
}

// Sets the horizontal adjustment of the viewport.
// Parameters:
// 	adjustment	a Adjustment.
func (v *CViewport) SetHAdjustment(adjustment *CAdjustment) {
	if err := v.SetStructProperty(PropertyHAdjustment, adjustment); err != nil {
		v.LogErr(err)
	}
}

// Sets the vertical adjustment of the viewport.
// Parameters:
// 	adjustment	a Adjustment.
func (v *CViewport) SetVAdjustment(adjustment *CAdjustment) {
	if err := v.SetStructProperty(PropertyVAdjustment, adjustment); err != nil {
		v.LogErr(err)
	}
}

// Sets the shadow type of the viewport.
// Parameters:
// 	type	the new shadow type.
func (v *CViewport) SetShadowType(shadowType ShadowType) {
	if err := v.SetStructProperty(PropertyViewportShadowType, shadowType); err != nil {
		v.LogErr(err)
	}
}

// Gets the shadow type of the Viewport. See
// SetShadowType.
// Returns:
// 	the shadow type
func (v *CViewport) GetShadowType() (shadowType ShadowType) {
	var ok bool
	if value, err := v.GetStructProperty(PropertyViewportShadowType); err != nil {
		v.LogErr(err)
	} else if shadowType, ok = value.(ShadowType); !ok {
		v.LogError("value stored in %v property is not of ShadowType type: %v (%T)", PropertyViewportShadowType, value, value)
	}
	return
}

// Gets the bin window of the Viewport.
// Returns:
// 	a Window.
// 	[transfer none]
func (v *CViewport) GetBinWindow() (value Window) {
	v.LogError("method unimplemented")
	return nil
}

// Gets the view window of the Viewport.
// Returns:
// 	a Window.
// 	[transfer none]
func (v *CViewport) GetViewWindow() (value Window) {
	v.LogError("method unimplemented")
	return nil
}

func (v *CViewport) Invalidate() cdk.EventFlag {
	if child := v.GetChild(); child != nil {
		local := child.GetOrigin()
		local.SubPoint(v.GetOrigin())
		if v.canvas == nil {
			v.canvas = cdk.NewCanvas(local, child.GetAllocation(), child.GetThemeRequest().Content.Normal)
		} else {
			v.canvas.SetOrigin(local)
			v.canvas.Resize(child.GetAllocation(), child.GetThemeRequest().Content.Normal)
		}
		return cdk.EVENT_STOP
	} else if v.canvas != nil {
		v.canvas = nil
	}
	return cdk.EVENT_PASS
}

func (v *CViewport) Resize() cdk.EventFlag {
	alloc := v.GetAllocation()
	child := v.GetChild()
	horizontal, vertical := v.GetHAdjustment(), v.GetVAdjustment()
	if alloc.W == 0 || alloc.H == 0 {
		if child != nil {
			child.SetAllocation(cdk.MakeRectangle(0, 0))
			child.Resize()
		}
		if horizontal != nil {
			horizontal.Configure(0, 0, 0, 0, 0, 0)
		}
		if vertical != nil {
			vertical.Configure(0, 0, 0, 0, 0, 0)
		}
		v.Invalidate()
		return v.Emit(SignalResize, v)
	}
	hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize := 0, 0, 0, 0, 0, 0
	vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize := 0, 0, 0, 0, 0, 0
	if child != nil {
		child.Freeze()
		defer child.Thaw()

		childSize := cdk.NewRectangle(child.GetSizeRequest())
		if childSize.W <= -1 {
			childSize.W = alloc.W
		}
		if childSize.H <= -1 {
			childSize.H = alloc.H
		}

		hChanged, vChanged := false, false
		if childSize.W >= alloc.W {
			delta := childSize.W - alloc.W
			hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize = 0, delta, 1, alloc.W/2, alloc.W
			if horizontal != nil {
				hValue = utils.ClampI(horizontal.GetValue(), 0, hUpper)
				horizontal.Configure(hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize)
			}
			hChanged = true
		}
		if childSize.H >= alloc.H {
			delta := childSize.H - alloc.H
			vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize = 0, delta, 1, alloc.H/2, alloc.H
			if vertical != nil {
				vValue = utils.ClampI(vertical.GetValue(), 0, vUpper)
				vertical.Configure(vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize)
			}
			vChanged = true
		}

		origin := v.GetOrigin()
		childOrigin := child.GetOrigin()
		if hChanged {
			childOrigin.X = origin.X - hValue
		}
		if vChanged {
			childOrigin.Y = origin.Y - vValue
		}

		child.SetOrigin(childOrigin.X, childOrigin.Y)
		child.SetAllocation(*childSize)

		if hChanged {
			if horizontal != nil {
				horizontal.Changed()
			}
		}
		if vChanged {
			if vertical != nil {
				vertical.Changed()
			}
		}
		v.Invalidate()
		return cdk.EVENT_STOP
	}
	return v.Emit(SignalResize, v)
}

func (v *CViewport) Draw(canvas cdk.Canvas) cdk.EventFlag {
	size := v.GetAllocation()
	if !v.IsVisible() || size.W <= 0 || size.H <= 0 {
		v.LogTrace("Label.Draw(): not visible, zero width or zero height")
		return cdk.EVENT_PASS
	}

	if child := v.GetChild(); child != nil {
		if f := child.Draw(v.canvas); f == cdk.EVENT_STOP {
			if err := canvas.Composite(v.canvas); err != nil {
				v.LogError("composite error: %v", err)
			}
		}

	}

	if debug, _ := v.GetBoolProperty(cdk.PropertyDebug); debug {
		canvas.DebugBox(cdk.ColorSilver, v.ObjectInfo())
	}
	return cdk.EVENT_STOP
}

// The Adjustment that determines the values of the horizontal position
// for this viewport.
// Flags: Read / Write / Construct
const PropertyViewportHAdjustment cdk.Property = "hadjustment"

// Determines how the shadowed box around the viewport is drawn.
// Flags: Read / Write
// Default value: GTK_SHADOW_IN
const PropertyViewportShadowType cdk.Property = "shadow-type"

// The Adjustment that determines the values of the vertical position for
// this viewport.
// Flags: Read / Write / Construct
const PropertyViewportVAdjustment cdk.Property = "vadjustment"

// Set the scroll adjustments for the viewport. Usually scrolled containers
// like ScrolledWindow will emit this signal to connect two instances of
// Scrollbar to the scroll directions of the Viewport.
// Listener function arguments:
// 	vertical Adjustment	the vertical GtkAdjustment
// 	arg2 Adjustment
const SignalSetScrollAdjustments cdk.Signal = "set-scroll-adjustments"
