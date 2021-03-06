package ctk

import (
	"fmt"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
)

const (
	TypeScrolledViewport cdk.CTypeTag = "ctk-scrolled-viewport"
)

func init() {
	_ = cdk.TypesManager.AddType(TypeScrolledViewport, func() interface{} { return MakeScrolledViewport() }, "ctk-scrolled-window")
	ctkBuilderTranslators[TypeScrolledViewport] = func(builder Builder, widget Widget, name, value string) error {
		switch name {
		case "hscrollbar-policy", "h-scrollbar-policy":
			if err := widget.SetPropertyFromString(PropertyHScrollbarPolicy, value); err != nil {
				return err
			}
			return nil
		case "vscrollbar-policy", "v-scrollbar-policy":
			if err := widget.SetPropertyFromString(PropertyVScrollbarPolicy, value); err != nil {
				return err
			}
			return nil
		}
		return ErrFallthrough
	}
}

// ScrolledViewport Hierarchy:
//      Object
//        +- Widget
//          +- Container
//            +- Bin
//              +- ScrolledViewport
type ScrolledViewport interface {
	Viewport

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	GetHAdjustment() (value *CAdjustment)
	GetVAdjustment() (value *CAdjustment)
	SetPolicy(hScrollbarPolicy PolicyType, vScrollbarPolicy PolicyType)
	AddWithViewport(child Widget)
	SetPlacement(windowPlacement CornerType)
	UnsetPlacement()
	SetShadowType(t ShadowType)
	SetHAdjustment(hAdjustment *CAdjustment)
	SetVAdjustment(vAdjustment *CAdjustment)
	GetPlacement() (value CornerType)
	GetPolicy() (hScrollbarPolicy PolicyType, vScrollbarPolicy PolicyType)
	GetShadowType() (value ShadowType)
	VerticalShowByPolicy() (show bool)
	HorizontalShowByPolicy() (show bool)
	SetTheme(theme cdk.Theme)
	Add(w Widget)
	Remove(w Widget)
	GetChild() Widget
	GetHScrollbar() *CHScrollbar
	GetVScrollbar() *CVScrollbar
	Show()
	Hide()
	GrabFocus()
	GetWidgetAt(p *cdk.Point2I) Widget
	InternalGetWidgetAt(p *cdk.Point2I) Widget
	CancelEvent()
	ProcessEvent(evt cdk.Event) cdk.EventFlag
	ProcessEventAtPoint(p *cdk.Point2I, evt *cdk.EventMouse) cdk.EventFlag
	GetRegions() (c, h, v cdk.Region)
	Resize() cdk.EventFlag
	Draw(canvas cdk.Canvas) cdk.EventFlag
	Invalidate() cdk.EventFlag
}

type CScrolledViewport struct {
	CViewport

	hCanvas               *cdk.CCanvas
	vCanvas               *cdk.CCanvas
	fCanvas               *cdk.CCanvas
	scrollbarsWithinBevel bool
	svFcHandle            string
}

func MakeScrolledViewport() *CScrolledViewport {
	return NewScrolledViewport()
}

func NewScrolledViewport() *CScrolledViewport {
	s := new(CScrolledViewport)
	s.Init()
	return s
}

func (s *CScrolledViewport) Init() (already bool) {
	if s.InitTypeItem(TypeScrolledViewport, s) {
		return true
	}
	s.CViewport.Init()
	s.flags = NULL_WIDGET_FLAG
	s.SetFlags(SENSITIVE)
	s.SetFlags(CAN_FOCUS)
	s.SetFlags(APP_PAINTABLE)
	_ = s.InstallProperty(PropertyHAdjustment, cdk.StructProperty, true, NewAdjustment(0, 0, 0, 0, 0, 0))
	_ = s.InstallProperty(PropertyHScrollbarPolicy, cdk.StructProperty, true, PolicyAlways)
	_ = s.InstallProperty(PropertyScrolledViewportShadowType, cdk.StructProperty, true, SHADOW_NONE)
	_ = s.InstallProperty(PropertyVAdjustment, cdk.StructProperty, true, NewAdjustment(0, 0, 0, 0, 0, 0))
	_ = s.InstallProperty(PropertyVScrollbarPolicy, cdk.StructProperty, true, PolicyAlways)
	_ = s.InstallProperty(PropertyWindowPlacement, cdk.StructProperty, true, GravityNorthWest)
	_ = s.InstallProperty(PropertyWindowPlacementSet, cdk.BoolProperty, true, false)
	s.CViewport.SetTheme(cdk.DefaultColorTheme)
	s.SetPolicy(PolicyAlways, PolicyAlways)
	// hScrollbar
	s.CContainer.Add(NewHScrollbar())
	s.hCanvas = cdk.NewCanvas(cdk.Point2I{}, cdk.Rectangle{}, cdk.DefaultColorTheme.Content.Normal)
	if hs := s.GetHScrollbar(); hs != nil {
		hs.SetParent(s)
		hs.SetWindow(s.GetWindow())
		s.SetHAdjustment(hs.GetAdjustment())
		hs.UnsetFlags(CAN_FOCUS)
	}
	// vScrollbar
	s.CContainer.Add(NewVScrollbar())
	s.vCanvas = cdk.NewCanvas(cdk.Point2I{}, cdk.Rectangle{}, cdk.DefaultColorTheme.Content.Normal)
	if vs := s.GetVScrollbar(); vs != nil {
		vs.SetParent(s)
		vs.SetWindow(s.GetWindow())
		s.SetVAdjustment(vs.GetAdjustment())
		vs.UnsetFlags(CAN_FOCUS)
	}
	s.svFcHandle = fmt.Sprintf("%v.focus-changed", s.ObjectName())
	s.Connect(SignalLostFocus, s.svFcHandle, s.handleLostFocus)
	s.Connect(SignalGainedFocus, s.svFcHandle, s.handleGainedFocus)
	s.Invalidate()
	return false
}

func (s *CScrolledViewport) Build(builder Builder, element *CBuilderElement) error {
	s.Freeze()
	defer s.Thaw()
	if err := s.CObject.Build(builder, element); err != nil {
		return err
	}
	if len(element.Children) > 0 {
		topChild := element.Children[0]
		if topClass, ok := topChild.Attributes["class"]; ok {
			switch topClass {
			case "GtkViewport":
				// GtkScrolledWindow -> GtkViewport -> Thing
				if len(topChild.Children) > 0 {
					grandchild := topChild.Children[0]
					if newChild := builder.Build(grandchild); newChild != nil {
						if grandWidget, ok := grandchild.Instance.(Widget); ok {
							s.Add(grandWidget)
						} else {
							s.LogError("viewport grandchild is not of Widget type: %v (%T)", grandchild, grandchild)
						}
					}
				} else {
					s.LogError("viewport child has no descendants")
				}
			default:
				// GtkScrolledWindow -> ScrollableThing
				if newChild := builder.Build(topChild); newChild != nil {
					if newWidget, ok := newChild.(Widget); ok {
						s.Add(newWidget)
					}
				}
			}
		}
	}
	return nil
}

// Returns the horizontal scrollbar's adjustment, used to connect the
// horizontal scrollbar to the child widget's horizontal scroll
// functionality.
// Returns:
//      the horizontal Adjustment.
//      [transfer none]
func (s *CScrolledViewport) GetHAdjustment() (value Adjustment) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyHAdjustment); err != nil {
		s.LogErr(err)
	} else if value, ok = v.(Adjustment); !ok {
		s.LogError("value stored in struct property is not of Adjustment type: %v (%T)", v, v)
	}
	return
}

// Returns the vertical scrollbar's adjustment, used to connect the vertical
// scrollbar to the child widget's vertical scroll functionality.
// Returns:
//      the vertical Adjustment.
//      [transfer none]
func (s *CScrolledViewport) GetVAdjustment() (value Adjustment) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyVAdjustment); err != nil {
		s.LogErr(err)
	} else if value, ok = v.(Adjustment); !ok {
		s.LogError("value stored in struct property is not an Adjustment: %v (%T)", v, v)
	}
	return
}

// Sets the scrollbar policy for the horizontal and vertical scrollbars. The
// policy determines when the scrollbar should appear; it is a value from the
// PolicyType enumeration. If GTK_POLICY_ALWAYS, the scrollbar is always
// present; if GTK_POLICY_NEVER, the scrollbar is never present; if
// GTK_POLICY_AUTOMATIC, the scrollbar is present only if needed (that is, if
// the slider part of the bar would be smaller than the trough - the display
// is larger than the page size).
// Parameters:
//      hScrollbarPolicy        policy for horizontal bar
//      vScrollbarPolicy        policy for vertical bar
func (s *CScrolledViewport) SetPolicy(hScrollbarPolicy PolicyType, vScrollbarPolicy PolicyType) {
	if err := s.SetStructProperty(PropertyHScrollbarPolicy, hScrollbarPolicy); err != nil {
		s.LogErr(err)
	}
	if err := s.SetStructProperty(PropertyVScrollbarPolicy, vScrollbarPolicy); err != nil {
		s.LogErr(err)
	}
	return
}

// Used to add children without native scrolling capabilities. This is simply
// a convenience function; it is equivalent to adding the unscrollable child
// to a viewport, then adding the viewport to the scrolled window. If a child
// has native scrolling, use ContainerAdd instead of this function.
// The viewport scrolls the child by moving its Window, and takes the size
// of the child to be the size of its toplevel Window. This will be very
// wrong for most widgets that support native scrolling; for example, if you
// add a widget such as TreeView with a viewport, the whole widget will
// scroll, including the column headings. Thus, widgets with native scrolling
// support should not be used with the Viewport proxy. A widget supports
// scrolling natively if the set_scroll_adjustments_signal field in
// WidgetClass is non-zero, i.e. has been filled in with a valid signal
// identifier.
// Parameters:
//      child   the widget you want to scroll
func (s *CScrolledViewport) AddWithViewport(child Widget) {}

// Sets the placement of the contents with respect to the scrollbars for the
// scrolled window. The default is GTK_CORNER_TOP_LEFT, meaning the child is
// in the top left, with the scrollbars underneath and to the right. Other
// values in CornerType are GTK_CORNER_TOP_RIGHT, GTK_CORNER_BOTTOM_LEFT,
// and GTK_CORNER_BOTTOM_RIGHT. See also GetPlacement
// and UnsetPlacement.
// Parameters:
//      windowPlacement position of the child window
func (s *CScrolledViewport) SetPlacement(windowPlacement CornerType) {}

// Unsets the placement of the contents with respect to the scrollbars for
// the scrolled window. If no window placement is set for a scrolled window,
// it obeys the "gtk-scrolled-window-placement" XSETTING. See also
// SetPlacement and
// GetPlacement.
func (s *CScrolledViewport) UnsetPlacement() {}

// Changes the type of shadow drawn around the contents of scrolled_window .
// Parameters:
//      type    kind of shadow to draw around scrolled window contents
func (s *CScrolledViewport) SetShadowType(t ShadowType) {
	if err := s.SetStructProperty(PropertyScrolledViewportShadowType, t); err != nil {
		s.LogErr(err)
	}
}

// Sets the Adjustment for the horizontal scrollbar.
// Parameters:
//      hAdjustment     horizontal scroll adjustment
func (s *CScrolledViewport) SetHAdjustment(hAdjustment *CAdjustment) {
	if err := s.SetStructProperty(PropertyHAdjustment, hAdjustment); err != nil {
		s.LogErr(err)
	}
}

// Sets the Adjustment for the vertical scrollbar.
// Parameters:
//      vAdjustment     vertical scroll adjustment
func (s *CScrolledViewport) SetVAdjustment(vAdjustment *CAdjustment) {
	if err := s.SetStructProperty(PropertyVAdjustment, vAdjustment); err != nil {
		s.LogErr(err)
	}
}

// Gets the placement of the contents with respect to the scrollbars for the
// scrolled window. See SetPlacement.
// Returns:
//      the current placement value.
//      See also SetPlacement and
//      UnsetPlacement.
func (s *CScrolledViewport) GetPlacement() (value CornerType) {
	return
}

// Retrieves the current policy values for the horizontal and vertical
// scrollbars. See SetPolicy.
// Parameters:
//      hScrollbarPolicy        location to store the policy
// for the horizontal scrollbar, or NULL.
//      vScrollbarPolicy        location to store the policy
// for the vertical scrollbar, or NULL.
func (s *CScrolledViewport) GetPolicy() (hScrollbarPolicy PolicyType, vScrollbarPolicy PolicyType) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyHScrollbarPolicy); err != nil {
		s.LogErr(err)
	} else if hScrollbarPolicy, ok = v.(PolicyType); !ok {
		s.LogError("value stored in struct property is not of PolicyType: %v (%T)", v, v)
	}
	ok = false
	if v, err := s.GetStructProperty(PropertyVScrollbarPolicy); err != nil {
		s.LogErr(err)
	} else if vScrollbarPolicy, ok = v.(PolicyType); !ok {
		s.LogError("value stored in struct property is not of PolicyType: %v (%T)", v, v)
	}
	return
}

// Gets the shadow type of the scrolled window. See
// SetShadowType.
// Returns:
//      the current shadow type
func (s *CScrolledViewport) GetShadowType() (value ShadowType) {
	var ok bool
	if v, err := s.GetStructProperty(PropertyScrolledViewportShadowType); err != nil {
		s.LogErr(err)
	} else if value, ok = v.(ShadowType); !ok {
		s.LogError("value stored in struct property is not of ShadowType: %v (%T)", v, v)
	}
	return
}

func (s *CScrolledViewport) VerticalShowByPolicy() (show bool) {
	vPolicy, _ := s.GetPolicy()
	if vertical := s.GetVAdjustment(); vertical != nil {
		show = vertical.ShowByPolicy(vPolicy)
		if !show && vPolicy == PolicyAutomatic && vertical.Moot() {
			if child := s.GetChild(); child != nil {
				childSize := cdk.NewRectangle(child.GetSizeRequest())
				if childSize.H > 0 {
					alloc := s.GetAllocation()
					if childSize.H > alloc.H {
						show = true
					}
				}
			}
		}
	} else {
		s.LogError("missing vertical adjustment")
	}
	return
}

func (s *CScrolledViewport) HorizontalShowByPolicy() (show bool) {
	_, hPolicy := s.GetPolicy()
	if horizontal := s.GetHAdjustment(); horizontal != nil {
		show = horizontal.ShowByPolicy(hPolicy)
		if !show && hPolicy == PolicyAutomatic && horizontal.Moot() {
			if child := s.GetChild(); child != nil {
				childSize := cdk.NewRectangle(child.GetSizeRequest())
				if childSize.W > 0 {
					alloc := s.GetAllocation()
					if childSize.W > alloc.W {
						show = true
					}
				}
			}
		}
	} else {
		s.LogError("missing horizontal adjustment")
	}
	return
}

func (s *CScrolledViewport) SetTheme(theme cdk.Theme) {
	s.CViewport.SetTheme(theme)
	s.Invalidate()
}

func (s *CScrolledViewport) Add(w Widget) {
	if len(s.children) < 3 {
		s.CContainer.Add(w)
		s.Invalidate()
	} else {
		s.LogError("too many children for scrolled viewport")
	}
}

func (s *CScrolledViewport) Remove(w Widget) {
	s.CContainer.Remove(w)
	s.Invalidate()
}

func (s *CScrolledViewport) GetChild() Widget {
	for _, child := range s.GetChildren() {
		if _, ok := child.(Scrollbar); !ok {
			return child
		}
	}
	return nil
}

func (s *CScrolledViewport) GetHScrollbar() *CHScrollbar {
	for _, child := range s.GetChildren() {
		if v, ok := child.(*CHScrollbar); ok {
			return v
		}
	}
	return nil
}

func (s *CScrolledViewport) GetVScrollbar() *CVScrollbar {
	for _, child := range s.GetChildren() {
		if v, ok := child.(*CVScrollbar); ok {
			return v
		}
	}
	return nil
}

func (s *CScrolledViewport) Show() {
	s.CViewport.Show()
	if child := s.GetChild(); child != nil {
		child.Show()
	}
	if hs := s.GetHScrollbar(); hs != nil {
		hs.Show()
	}
	if vs := s.GetVScrollbar(); vs != nil {
		vs.Show()
	}
	s.Invalidate()
}

func (s *CScrolledViewport) Hide() {
	s.CViewport.Hide()
	if child := s.GetChild(); child != nil {
		child.Hide()
	}
	if hs := s.GetHScrollbar(); hs != nil {
		hs.Hide()
	}
	if vs := s.GetVScrollbar(); vs != nil {
		vs.Hide()
	}
	s.Invalidate()
}

// If the Widget instance CanFocus() then take the focus of the associated
// Window. Any previously focused Widget will emit a lost-focus signal and the
// newly focused Widget will emit a gained-focus signal. This method emits a
// grab-focus signal initially and if the listeners return EVENT_PASS, the
// changes are applied
//
// Emits: SignalGrabFocus, Argv=[Widget instance]
// Emits: SignalLostFocus, Argv=[Previous focus Widget instance], From=Previous focus Widget instance
// Emits: SignalGainedFocus, Argv=[Widget instance, previous focus Widget instance]
func (s *CScrolledViewport) GrabFocus() {
	if s.CanFocus() {
		if r := s.Emit(SignalGrabFocus, s); r == cdk.EVENT_PASS {
			tl := s.GetWindow()
			if tl != nil {
				var fw Widget
				focused := tl.GetFocus()
				tl.SetFocus(s)
				if focused != nil {
					var ok bool
					if fw, ok = focused.(Widget); ok && fw.ObjectID() != s.ObjectID() {
						if f := fw.Emit(SignalLostFocus, fw); f == cdk.EVENT_STOP {
							fw = nil
						}
					} else if ok {
						// already the focus, nothing to do
						return
					}
				}
				if f := s.Emit(SignalGainedFocus, s, fw); f == cdk.EVENT_STOP {
					if fw != nil {
						tl.SetFocus(fw)
					}
				}
			}
		}
	}
}

func (s *CScrolledViewport) GetWidgetAt(p *cdk.Point2I) Widget {
	if s.HasPoint(p) && s.IsVisible() {
		return s
	}
	return nil
}

func (s *CScrolledViewport) InternalGetWidgetAt(p *cdk.Point2I) Widget {
	if s.HasPoint(p) {
		if vs := s.GetVScrollbar(); vs != nil {
			if vs.HasPoint(p) {
				return vs
			}
		}
		if hs := s.GetHScrollbar(); hs != nil {
			if hs.HasPoint(p) {
				return hs
			}
		}
		if child := s.GetChild(); child != nil {
			if child.HasPoint(p) {
				return child
			}
		}
		return s
	}
	return nil
}

func (s *CScrolledViewport) CancelEvent() {
	if child := s.GetChild(); child != nil {
		if cs, ok := child.(Sensitive); ok {
			cs.CancelEvent()
		}
	}
	if vs := s.GetVScrollbar(); vs != nil {
		vs.CancelEvent()
	}
	if hs := s.GetHScrollbar(); hs != nil {
		hs.CancelEvent()
	}
	s.Invalidate()
}

func (s *CScrolledViewport) ProcessEvent(evt cdk.Event) cdk.EventFlag {
	s.Lock()
	defer s.Unlock()
	switch e := evt.(type) {
	case *cdk.EventMouse:
		if e.IsWheelImpulse() {
			s.GrabFocus()
			switch e.WheelImpulse() {
			case cdk.WheelUp:
				if vs := s.GetVScrollbar(); vs != nil {
					if f := vs.ForwardStep(); f == cdk.EVENT_STOP {
						s.Invalidate()
						return cdk.EVENT_STOP
					}
					return cdk.EVENT_PASS
				}
			case cdk.WheelLeft:
				if hs := s.GetHScrollbar(); hs != nil {
					if f := hs.BackwardStep(); f == cdk.EVENT_STOP {
						s.Invalidate()
						return cdk.EVENT_STOP
					}
					return cdk.EVENT_PASS
				}
			case cdk.WheelDown:
				if vs := s.GetVScrollbar(); vs != nil {
					if f := vs.BackwardStep(); f == cdk.EVENT_STOP {
						s.Invalidate()
						return cdk.EVENT_STOP
					}
				}
			case cdk.WheelRight:
				if hs := s.GetHScrollbar(); hs != nil {
					if f := hs.ForwardStep(); f == cdk.EVENT_STOP {
						s.Invalidate()
						return cdk.EVENT_STOP
					}
					return cdk.EVENT_PASS
				}
			}
		}
		point := cdk.NewPoint2I(e.Position())
		if f := s.ProcessEventAtPoint(point, e); f == cdk.EVENT_STOP {
			s.GrabFocus()
			return cdk.EVENT_STOP
		}
	case *cdk.EventKey:
		if vs := s.GetVScrollbar(); vs != nil {
			if f := vs.ProcessEvent(evt); f == cdk.EVENT_STOP {
				s.Invalidate()
				return cdk.EVENT_STOP
			}
		}
		if hs := s.GetHScrollbar(); hs != nil {
			if f := hs.ProcessEvent(evt); f == cdk.EVENT_STOP {
				s.Invalidate()
				return cdk.EVENT_STOP
			}
		}
	}
	return cdk.EVENT_PASS
}

func (s *CScrolledViewport) ProcessEventAtPoint(p *cdk.Point2I, evt *cdk.EventMouse) cdk.EventFlag {
	if w := s.InternalGetWidgetAt(p); w != nil {
		if w.ObjectID() != s.ObjectID() {
			if ws, ok := w.(Sensitive); ok {
				if f := ws.ProcessEvent(evt); f == cdk.EVENT_STOP {
					s.Invalidate()
					return cdk.EVENT_STOP
				}
			}
		}
	}
	return cdk.EVENT_PASS
}

//
// func (s *CScrolledViewport) SetSizeRequest(w, h int) {
// 	s.CWidget.SetSizeRequest(w, h)
// 	s.Invalidate()
// }

// Returns a CDK Region for each of the viewport child space, horizontal and vertical
// scrollbar spaces.
func (s *CScrolledViewport) GetRegions() (c, h, v cdk.Region) {
	if child := s.GetChild(); child != nil {
		o := child.GetOrigin()
		a := child.GetAllocation()
		c = cdk.MakeRegion(o.X, o.Y, a.W, a.H)
	}
	if hs := s.GetHScrollbar(); hs != nil && s.HorizontalShowByPolicy() {
		o := hs.GetOrigin()
		a := hs.GetAllocation()
		h = cdk.MakeRegion(o.X, o.Y, a.W, a.H)
	}
	if vs := s.GetVScrollbar(); vs != nil && s.VerticalShowByPolicy() {
		o := vs.GetOrigin()
		a := vs.GetAllocation()
		v = cdk.MakeRegion(o.X, o.Y, a.W, a.H)
	}
	return
}

func (s *CScrolledViewport) Resize() cdk.EventFlag {
	// s.resizeViewport()
	// s.resizeScrollbars()
	s.Invalidate()
	return cdk.EVENT_STOP
}

func (s *CScrolledViewport) Draw(canvas cdk.Canvas) cdk.EventFlag {
	s.Lock()
	defer s.Unlock()
	alloc := s.GetAllocation()
	if !s.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
		return cdk.EVENT_PASS
	}
	child := s.GetChild()
	if child != nil {
		canvas.BoxWithTheme(
			s.GetOrigin(),
			s.GetAllocation(),
			false,
			true,
			child.GetTheme(),
		)
		child.Draw(s.canvas)
		if err := s.fCanvas.Composite(s.canvas); err != nil {
			s.LogError("child composite error: %v", err)
		}
		if err := canvas.Composite(s.fCanvas); err != nil {
			s.LogError("viewport composite error: %v", err)
		}
	}
	if vs := s.GetVScrollbar(); child != nil && vs != nil && s.VerticalShowByPolicy() {
		vs.Draw(s.vCanvas)
		if err := canvas.Composite(s.vCanvas); err != nil {
			s.LogError("vertical scrollbar composite error: %v", err)
		}
	}
	if hs := s.GetHScrollbar(); child != nil && hs != nil && s.HorizontalShowByPolicy() {
		hs.Draw(s.hCanvas)
		if err := canvas.Composite(s.hCanvas); err != nil {
			s.LogError("horizontal scrollbar composite error: %v", err)
		}
	}
	if child != nil && s.VerticalShowByPolicy() && s.HorizontalShowByPolicy() {
		// fill in the corner gap between scrollbars
		_ = canvas.SetRune(alloc.W-1, alloc.H-1, s.GetTheme().Content.FillRune, s.GetTheme().Content.Normal)
	}

	if debug, _ := s.GetBoolProperty(cdk.PropertyDebug); debug {
		canvas.DebugBox(cdk.ColorSilver, s.ObjectInfo())
	}
	return cdk.EVENT_STOP
}

func (s *CScrolledViewport) Invalidate() cdk.EventFlag {
	s.resizeViewport()
	s.resizeScrollbars()
	alloc := s.GetAllocation()
	origin := s.GetOrigin()
	if child := s.GetChild(); child != nil {
		local := child.GetOrigin()
		local.SubPoint(origin)
		size := child.GetAllocation() // set by resizeViewport() call
		if s.canvas == nil {
			s.canvas = cdk.NewCanvas(local, size, child.GetTheme().Content.Normal)
		} else {
			s.canvas.SetOrigin(local)
			s.canvas.Resize(size, child.GetTheme().Content.Normal)
		}
		size = alloc.Clone()
		size.Clamp(0, 0, alloc.W, alloc.H)
		if s.HorizontalShowByPolicy() {
			size.H -= 1
		}
		if s.VerticalShowByPolicy() {
			size.W -= 1
		}
		local = cdk.MakePoint2I(0, 0)
		if s.fCanvas == nil {
			s.fCanvas = cdk.NewCanvas(local, size, child.GetTheme().Content.Normal)
		} else {
			s.fCanvas.SetOrigin(local)
			s.fCanvas.Resize(size, child.GetTheme().Content.Normal)
		}
	}
	if vs := s.GetVScrollbar(); vs != nil && s.VerticalShowByPolicy() {
		local := vs.GetOrigin()
		local.SubPoint(origin)
		s.vCanvas.SetOrigin(local)
		s.vCanvas.Resize(vs.GetAllocation(), s.GetTheme().Content.Normal)
		vs.Show()
	}
	if hs := s.GetHScrollbar(); hs != nil && s.HorizontalShowByPolicy() {
		local := hs.GetOrigin()
		local.SubPoint(origin)
		s.hCanvas.SetOrigin(local)
		s.hCanvas.Resize(hs.GetAllocation(), s.GetTheme().Content.Normal)
		hs.Show()
	}
	return cdk.EVENT_STOP
}

func (s *CScrolledViewport) makeAdjustments() (region cdk.Region, changed bool) {
	changed = false
	region = cdk.MakeRegion(0, 0, 0, 0)
	origin := s.GetOrigin()
	alloc := s.GetAllocation()
	horizontal, vertical := s.GetHAdjustment(), s.GetVAdjustment()
	if alloc.W == 0 || alloc.H == 0 {
		if horizontal != nil {
			ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize := horizontal.Settings()
			ah := []int{ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize}
			bh := []int{0, 0, 0, 0, 0, 0}
			changed = utils.EqInts(ah, bh)
			horizontal.Configure(0, 0, 0, 0, 0, 0)
		}
		if vertical != nil {
			ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize := vertical.Settings()
			av := []int{ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize}
			bv := []int{0, 0, 0, 0, 0, 0}
			changed = changed || utils.EqInts(av, bv)
			vertical.Configure(0, 0, 0, 0, 0, 0)
		}
		return
	}
	hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize := 0, 0, 0, 0, 0, 0
	vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize := 0, 0, 0, 0, 0, 0
	if child := s.GetChild(); child != nil {
		size := cdk.NewRectangle(child.GetSizeRequest())
		if size.W <= -1 { // auto
			size.W = alloc.W
			if s.VerticalShowByPolicy() {
				size.W -= 1
			}
		}
		if size.H <= -1 { // auto
			size.H = alloc.H
			if s.HorizontalShowByPolicy() {
				size.H -= 1
			}
		}
		if size.W >= alloc.W {
			hStepIncrement, hPageIncrement, hPageSize = 1, alloc.W/2, alloc.W
			if size.W >= alloc.W {
				overflow := size.W - alloc.W
				hLower, hUpper = 0, overflow
				if s.VerticalShowByPolicy() {
					hUpper += 1
				}
			} else {
				hLower, hUpper, hValue = 0, 0, 0
			}
			if horizontal != nil {
				hValue = utils.ClampI(horizontal.GetValue(), hLower, hUpper)
			}
		}
		region.X = origin.X - hValue
		region.W = size.W
		if size.H >= alloc.H {
			vStepIncrement, vPageIncrement, vPageSize = 1, alloc.H/2, alloc.H
			if size.H >= alloc.H {
				overflow := size.H - alloc.H
				vLower, vUpper = 0, overflow
				if s.HorizontalShowByPolicy() {
					vUpper += 1
				}
			} else {
				vLower, vUpper, vValue = 0, 0, 0
			}
			if vertical != nil {
				vValue = utils.ClampI(vertical.GetValue(), vLower, vUpper)
			}
		}
		region.Y = origin.Y - vValue
		region.H = size.H
	}
	// horizontal
	if horizontal != nil {
		ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize := horizontal.Settings()
		ah := []int{ohValue, ohLower, ohUpper, ohStepIncrement, ohPageIncrement, ohPageSize}
		bh := []int{hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize}
		if !utils.EqInts(ah, bh) {
			changed = true
			horizontal.Configure(hValue, hLower, hUpper, hStepIncrement, hPageIncrement, hPageSize)
		}
	}
	// vertical
	if vertical != nil {
		ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize := vertical.Settings()
		av := []int{ovValue, ovLower, ovUpper, ovStepIncrement, ovPageIncrement, ovPageSize}
		bv := []int{vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize}
		if !utils.EqInts(av, bv) {
			changed = true
			vertical.Configure(vValue, vLower, vUpper, vStepIncrement, vPageIncrement, vPageSize)
		}
	}
	return
}

func (s *CScrolledViewport) resizeViewport() cdk.EventFlag {
	region, _ := s.makeAdjustments()
	if child := s.GetChild(); child != nil {
		child.SetOrigin(region.X, region.Y)
		child.SetAllocation(region.Size())
		return child.Resize()
	}
	return cdk.EVENT_STOP
}

func (s *CScrolledViewport) resizeScrollbars() cdk.EventFlag {
	origin := s.GetOrigin()
	alloc := s.GetAllocation()
	if hs := s.GetHScrollbar(); hs != nil {
		o := cdk.MakePoint2I(origin.X, origin.Y+alloc.H-1)
		a := cdk.MakeRectangle(alloc.W, 1)
		if s.VerticalShowByPolicy() {
			a.W -= 1
		}
		hs.SetOrigin(o.X, o.Y)
		hs.SetAllocation(a)
		theme := hs.GetTheme()
		if s.IsFocused() {
			theme.Content.Normal = theme.Content.Focused
			theme.Border.Normal = theme.Border.Focused
		}
		hs.SetThemeRequest(theme)
		hs.Resize()
	}
	if vs := s.GetVScrollbar(); vs != nil {
		o := cdk.MakePoint2I(origin.X+alloc.W-1, origin.Y)
		a := cdk.MakeRectangle(1, alloc.H)
		if s.HorizontalShowByPolicy() {
			a.H -= 1
		}
		vs.SetOrigin(o.X, o.Y)
		vs.SetAllocation(a)
		theme := vs.GetTheme()
		if s.IsFocused() {
			theme.Content.Normal = theme.Content.Focused
			theme.Border.Normal = theme.Border.Focused
		}
		vs.SetThemeRequest(theme)
		vs.Resize()
	}
	return cdk.EVENT_STOP
}

func (s *CScrolledViewport) handleLostFocus([]interface{}, ...interface{}) cdk.EventFlag {
	s.Invalidate()
	return cdk.EVENT_PASS
}

func (s *CScrolledViewport) handleGainedFocus([]interface{}, ...interface{}) cdk.EventFlag {
	s.Invalidate()
	return cdk.EVENT_PASS
}

// The Adjustment for the horizontal position.
// Flags: Read / Write / Construct
const PropertyHAdjustment cdk.Property = "h-adjustment"

// When the horizontal scrollbar is displayed.
// Flags: Read / Write
// Default value: GTK_POLICY_ALWAYS
const PropertyHScrollbarPolicy cdk.Property = "h-scrollbar-policy"

// Style of bevel around the contents.
// Flags: Read / Write
// Default value: GTK_SHADOW_NONE
const PropertyScrolledViewportShadowType cdk.Property = "viewport-shadow-type"

// The Adjustment for the vertical position.
// Flags: Read / Write / Construct
const PropertyVAdjustment cdk.Property = "v-adjustment"

// When the vertical scrollbar is displayed.
// Flags: Read / Write
// Default value: GTK_POLICY_ALWAYS
const PropertyVScrollbarPolicy cdk.Property = "vscrollbar-policy"

// Where the contents are located with respect to the scrollbars. This
// property only takes effect if "window-placement-set" is TRUE.
// Flags: Read / Write
// Default value: GTK_CORNER_TOP_LEFT
const PropertyWindowPlacement cdk.Property = "window-placement"

// Whether "window-placement" should be used to determine the location of the
// contents with respect to the scrollbars. Otherwise, the
// "gtk-scrolled-window-placement" setting is used.
// Flags: Read / Write
// Default value: FALSE
const PropertyWindowPlacementSet cdk.Property = "window-placement-set"

// Listener function arguments:
//      arg1 DirectionType
const SignalMoveFocusOut cdk.Signal = "move-focus-out"

// The ::scroll-child signal is a which gets emitted when a keybinding that
// scrolls is pressed. The horizontal or vertical adjustment is updated which
// triggers a signal that the scrolled windows child may listen to and scroll
// itself.
const SignalScrollChild cdk.Signal = "scroll-child"
