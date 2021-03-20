package ctk

import (
	"fmt"

	"github.com/kckrinke/go-cdk"
)

// CDK type-tag for Frame objects
const TypeFrame cdk.CTypeTag = "ctk-frame"

func init() {
	_ = cdk.TypesManager.AddType(TypeFrame, func() interface{} { return MakeFrame() })
}

// Frame Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Frame
//	          +- AspectFrame
type Frame interface {
	Bin
	Buildable

	Init() (already bool)
	SetLabel(label string)
	SetLabelWidget(labelWidget Widget)
	SetLabelAlign(xAlign float64, yAlign float64)
	SetShadowType(shadowType ShadowType)
	GetLabel() (value string)
	GetLabelAlign() (xAlign float64, yAlign float64)
	GetLabelWidget() (value Widget)
	GetShadowType() (value ShadowType)
	GrabFocus()
	Add(w Widget)
	Remove(w Widget)
	IsFocus() bool
	GetFocusWithChild() (focusWithChild bool)
	SetFocusWithChild(focusWithChild bool)
	GetWidgetAt(p *cdk.Point2I) Widget
	GetThemeRequest() (theme cdk.Theme)
	GetSizeRequest() (width, height int)
	Resize() cdk.EventFlag
	Draw(canvas cdk.Canvas) cdk.EventFlag
	Invalidate() cdk.EventFlag
}

// The CFrame structure implements the Frame interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Frame objects
type CFrame struct {
	CBin

	labelCanvas    *cdk.CCanvas
	canvas         *cdk.CCanvas
	focusWithChild bool
	fFcHandle      string
}

// Default constructor for Frame objects
func MakeFrame() *CFrame {
	return NewFrame("")
}

// construct a new frame with a default widget containing the given text
func NewFrame(text string) *CFrame {
	f := new(CFrame)
	f.Init()
	label := NewLabel(text)
	label.SetSingleLineMode(true)
	label.SetLineWrap(false)
	label.SetLineWrapMode(cdk.WRAP_NONE)
	label.SetJustify(cdk.JUSTIFY_LEFT)
	label.SetParent(f)
	label.SetWindow(f.GetWindow())
	label.Show()
	f.SetLabelWidget(label)
	return f
}

// construct a new frame with the given widget instead of the normal label
func NewFrameWithWidget(w Widget) *CFrame {
	f := new(CFrame)
	f.Init()
	f.SetLabelWidget(w)
	return f
}

// Frame object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Frame instance
func (f *CFrame) Init() (already bool) {
	if f.InitTypeItem(TypeFrame, f) {
		return true
	}
	f.CBin.Init()
	f.flags = NULL_WIDGET_FLAG
	f.SetFlags(PARENT_SENSITIVE)
	f.SetFlags(APP_PAINTABLE)
	_ = f.InstallProperty(PropertyLabel, cdk.StringProperty, true, nil)
	_ = f.InstallProperty(PropertyLabelWidget, cdk.StructProperty, true, nil)
	_ = f.InstallProperty(PropertyLabelXAlign, cdk.FloatProperty, true, 0.0)
	_ = f.InstallProperty(PropertyLabelYAlign, cdk.FloatProperty, true, 0.5)
	_ = f.InstallProperty(PropertyShadow, cdk.StructProperty, true, nil)
	_ = f.InstallProperty(PropertyShadowType, cdk.StructProperty, true, nil)
	f.focusWithChild = false
	f.labelCanvas = cdk.NewCanvas(cdk.Point2I{}, cdk.Rectangle{}, f.GetTheme().Content.Normal)
	f.fFcHandle = fmt.Sprintf("%v-frame.focus-changed", f.ObjectName())
	return false
}

// Sets the text of the label. If label is NULL, the current label is
// removed.
// Parameters:
// 	label	the text to use as the label of the frame.
func (f *CFrame) SetLabel(label string) {
	if err := f.SetStringProperty(PropertyLabel, label); err != nil {
		f.LogErr(err)
	} else {
		if w := f.GetLabelWidget(); w != nil {
			if lw, ok := w.(Label); ok {
				lw.SetText(label)
				f.Invalidate()
			}
		}
	}
}

// Sets the label widget for the frame. This is the widget that will appear
// embedded in the top edge of the frame as a title.
// Parameters:
// 	labelWidget	the new label widget
func (f *CFrame) SetLabelWidget(labelWidget Widget) {
	if err := f.SetStructProperty(PropertyLabelWidget, labelWidget); err != nil {
		f.LogErr(err)
	} else {
		labelWidget.SetParent(f)
		labelWidget.SetWindow(f.GetWindow())
		labelWidget.Show()
		f.Invalidate()
	}
}

// Sets the alignment of the frame widget's label. The default values for a
// newly created frame are 0.0 and 0.5.
// Parameters:
// 	xAlign	The position of the label along the top edge of the widget. A value
// 	        of 0.0 represents left alignment; 1.0 represents right alignment.
// 	yAlign	The y alignment of the label. A value of 0.0 aligns under the frame;
// 	        1.0 aligns above the frame. If the values are exactly 0.0 or 1.0 the
// 	        gap in the frame won't be painted because the label will be
// 	        completely above or below the frame.
func (f *CFrame) SetLabelAlign(xAlign float64, yAlign float64) {
	if err := f.SetFloatProperty(PropertyLabelXAlign, xAlign); err != nil {
		f.LogErr(err)
	}
	if err := f.SetFloatProperty(PropertyLabelYAlign, yAlign); err != nil {
		f.LogErr(err)
	}
}

// Sets the shadow type for frame .
// Parameters:
// 	type	the new ShadowType
func (f *CFrame) SetShadowType(shadowType ShadowType) {
	if err := f.SetStructProperty(PropertyShadowType, shadowType); err != nil {
		f.LogErr(err)
	}
}

// If the frame's label widget is a Label, returns the text in the label
// widget. (The frame will have a Label for the label widget if a non-NULL
// argument was passed to New.)
// Returns:
// 	the text in the label, or NULL if there was no label widget or
// 	the lable widget was not a Label. This string is owned by
// 	CTK and must not be modified or freed.
func (f *CFrame) GetLabel() (value string) {
	var err error
	if value, err = f.GetStringProperty(PropertyLabel); err != nil {
		f.LogErr(err)
	}
	return
}

// Retrieves the X and Y alignment of the frame's label. See SetLabelAlign.
// Parameters:
// 	xAlign	X alignment of frame's label
// 	yAlign	Y alignment of frame's label
func (f *CFrame) GetLabelAlign() (xAlign float64, yAlign float64) {
	var err error
	if xAlign, err = f.GetFloatProperty(PropertyLabelXAlign); err != nil {
		f.LogErr(err)
	}
	if yAlign, err = f.GetFloatProperty(PropertyLabelYAlign); err != nil {
		f.LogErr(err)
	}
	return
}

// Retrieves the label widget for the frame. See
// SetLabelWidget.
// Returns:
// 	the label widget, or NULL if there is none.
// 	[transfer none]
func (f *CFrame) GetLabelWidget() (value Widget) {
	if v, err := f.GetStructProperty(PropertyLabelWidget); err != nil {
		f.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(Widget); !ok {
			f.LogError("value stored in %v is not a Widget: %v (%T)", PropertyLabelWidget, v, v)
		}
	}
	return
}

// Retrieves the shadow type of the frame. See SetShadowType.
// Returns:
// 	the current shadow type of the frame.
func (f *CFrame) GetShadowType() (value ShadowType) {
	if v, err := f.GetStructProperty(PropertyShadowType); err != nil {
		f.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(ShadowType); !ok {
			f.LogError("value stored in %v is not a ShadowType: %v (%T)", PropertyShadowType, v, v)
		}
	}
	return
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
func (f *CFrame) GrabFocus() {
	if f.CanFocus() {
		if r := f.Emit(SignalGrabFocus, f); r == cdk.EVENT_PASS {
			tl := f.GetWindow()
			if tl != nil {
				var fw Widget
				focused := tl.GetFocus()
				tl.SetFocus(f)
				if focused != nil {
					var ok bool
					if fw, ok = focused.(Widget); ok && fw.ObjectID() != f.ObjectID() {
						if f := fw.Emit(SignalLostFocus, fw); f == cdk.EVENT_STOP {
							fw = nil
						}
					}
				}
				if f := f.Emit(SignalGainedFocus, f, fw); f == cdk.EVENT_STOP {
					if fw != nil {
						tl.SetFocus(fw)
					}
				}
				f.LogDebug("has taken focus")
			}
		}
	}
}

func (f *CFrame) Add(w Widget) {
	f.CBin.Add(w)
	w.Connect(SignalLostFocus, f.fFcHandle, f.handleLostFocus)
	w.Connect(SignalGainedFocus, f.fFcHandle, f.handleGainedFocus)
	f.Invalidate()
}

func (f *CFrame) Remove(w Widget) {
	_ = w.Disconnect(SignalLostFocus, f.fFcHandle)
	_ = w.Disconnect(SignalGainedFocus, f.fFcHandle)
	f.CBin.Remove(w)
	f.Invalidate()
}

func (f *CFrame) IsFocus() bool {
	if child := f.GetChild(); child != nil && child.CanFocus() {
		return child.IsFocus()
	}
	return f.CBin.IsFocus()
}

func (f *CFrame) GetFocusWithChild() (focusWithChild bool) {
	return f.focusWithChild
}

func (f *CFrame) SetFocusWithChild(focusWithChild bool) {
	f.focusWithChild = focusWithChild
	f.Invalidate()
}

func (f *CFrame) GetWidgetAt(p *cdk.Point2I) Widget {
	if f.HasPoint(p) && f.IsVisible() {
		if child := f.GetChild(); child != nil {
			if cc, ok := child.(Container); ok {
				if cc.HasPoint(p) && cc.IsVisible() {
					if w := cc.GetWidgetAt(p); w != nil && w.IsVisible() {
						return w
					}
				}
			}
		}
		return f
	}
	return nil
}

// Returns the current theme, adjusted for Widget focus and accounting for
// any PARENT_SENSITIVE conditions. This method is primarily useful in drawable
// Widget types during the Invalidate() and Draw() stages of the Widget
// lifecycle
func (f *CFrame) GetThemeRequest() (theme cdk.Theme) {
	theme = f.GetTheme()
	var childFocused bool
	if child := f.GetChild(); child != nil {
		childFocused = child.IsFocus()
	}
	if childFocused || (f.CanFocus() && f.IsFocused()) || f.IsParentFocused() {
		theme.Content.Normal = theme.Content.Focused
		theme.Border.Normal = theme.Border.Focused
	}
	return
}

func (f *CFrame) GetSizeRequest() (width, height int) {
	_, yAlign := f.GetLabelAlign()
	size := cdk.NewRectangle(f.CWidget.GetSizeRequest())
	if child := f.GetChild(); child != nil {
		childSize := cdk.NewRectangle(child.GetSizeRequest())
		if size.W <= -1 {
			if childSize.W > -1 {
				size.W = 1 + childSize.W + 1
			}
		}
		if size.H <= -1 {
			if childSize.H > -1 {
				size.H = 1 + childSize.H + 1
				if yAlign == 0.0 {
					size.H += 1
				}
			}
		}
	}
	return size.W, size.H
}

func (f *CFrame) Resize() cdk.EventFlag {
	f.Lock()
	defer f.Unlock()
	// our allocation has been set prior to Resize() being called
	alloc := f.GetAllocation()
	widget := f.GetLabelWidget()
	if widget == nil {
		return cdk.EVENT_PASS
	}
	label, _ := widget.(Label)
	if alloc.W <= 0 && alloc.H <= 0 {
		if label != nil {
			label.SetAllocation(cdk.MakeRectangle(0, 0))
			label.Resize()
		}
		return cdk.EVENT_PASS
	}
	alloc.Floor(0, 0)
	origin := f.GetOrigin()
	childOrigin := cdk.MakePoint2I(origin.X+1, origin.Y+1)
	childAlloc := cdk.MakeRectangle(alloc.W-2, alloc.H-2)
	labelOrigin := cdk.MakePoint2I(origin.X+2, origin.Y)
	labelAlloc := cdk.MakeRectangle(alloc.W-4, 1)
	xAlign, yAlign := f.GetLabelAlign()
	if yAlign <= 0.0 {
		yAlign = 0.0
		childAlloc.H -= 1
		childOrigin.Y += 1
	} else if yAlign >= 1.0 {
		yAlign = 1.0
		labelOrigin.Y += 1
		childAlloc.H -= 1
		childOrigin.Y += 1
	} else {
		yAlign = 0.5
	}
	if label != nil {
		label.SetAlignment(xAlign, yAlign)
		label.SetMaxWidthChars(labelAlloc.W)
		label.SetOrigin(labelOrigin.X, labelOrigin.Y)
		label.SetAllocation(labelAlloc)
		label.Resize()
	}
	if child := f.GetChild(); child != nil {
		child.SetOrigin(childOrigin.X, childOrigin.Y)
		child.SetAllocation(childAlloc)
		child.Resize()
	}
	return f.Invalidate()
}

func (f *CFrame) Draw(canvas cdk.Canvas) cdk.EventFlag {
	f.Lock()
	defer f.Unlock()
	alloc := f.GetAllocation()
	if !f.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
		f.LogTrace("Label.Draw(): not visible, zero width or zero height")
		return cdk.EVENT_PASS
	}

	// render the box and border, with widget
	child := f.GetChild()
	theme := f.GetThemeRequest()
	if child != nil {
		theme = child.GetThemeRequest()
	}
	boxOrigin := cdk.MakePoint2I(0, 0)
	boxSize := alloc
	_, yAlign := f.GetLabelAlign()
	if yAlign == 0.0 { // top
		boxOrigin.Y += 1
		boxSize.H -= 1
	}
	canvas.BoxWithTheme(boxOrigin, boxSize, true, true, theme)

	if widget := f.GetLabelWidget(); widget != nil {
		if label, ok := widget.(Label); ok {
			labelTheme := label.GetTheme()
			if labelTheme.String() != theme.String() {
				label.SetTheme(theme)
				label.Invalidate()
			}
			label.Draw(f.labelCanvas)
			if err := canvas.Composite(f.labelCanvas); err != nil {
				f.LogError("composite error: %v", err)
			}
		}
	}

	if child != nil {
		child.Draw(f.canvas)
		if err := canvas.Composite(f.canvas); err != nil {
			f.LogError("composite error: %v", err)
		}
	}

	if debug, _ := f.GetBoolProperty(cdk.PropertyDebug); debug {
		canvas.DebugBox(cdk.ColorSilver, f.ObjectInfo())
	}
	return cdk.EVENT_STOP
}

func (f *CFrame) Invalidate() cdk.EventFlag {
	wantStop := false
	origin := f.GetOrigin()
	theme := f.GetThemeRequest()
	if labelChild, ok := f.GetLabelWidget().(Label); ok && labelChild != nil {
		local := labelChild.GetOrigin()
		local.SubPoint(origin)
		alloc := labelChild.GetAllocation()
		if f.labelCanvas == nil {
			f.labelCanvas = cdk.NewCanvas(local, alloc, theme.Content.Normal)
		} else {
			f.labelCanvas.SetOrigin(local)
			f.labelCanvas.Resize(alloc, theme.Content.Normal)
		}
		theme.Content.FillRune = rune(0)
		f.labelCanvas.Fill(theme)
		labelChild.SetTheme(theme)
		labelChild.Invalidate()
		wantStop = true
	}
	if child := f.GetChild(); child != nil {
		local := child.GetOrigin()
		local.SubPoint(origin)
		alloc := child.GetAllocation()
		if f.canvas == nil {
			f.canvas = cdk.NewCanvas(local, alloc, theme.Content.Normal)
		} else {
			f.canvas.SetOrigin(local)
			f.canvas.Resize(alloc, theme.Content.Normal)
		}
		wantStop = true
	}
	if wantStop {
		return cdk.EVENT_STOP
	}
	return cdk.EVENT_PASS
}

func (f *CFrame) handleLostFocus(_ []interface{}, _ ...interface{}) cdk.EventFlag {
	// f.Lock()
	// defer f.Unlock()
	theme := f.GetTheme()
	if f.focusWithChild {
		if child := f.GetChild(); child != nil && child.IsFocus() {
			theme.Content.Normal = theme.Content.Focused
			theme.Border.Normal = theme.Border.Focused
		}
	} else {
		// theme defaults to whatever is set
	}
	f.SetThemeRequest(theme)
	f.Invalidate()
	return cdk.EVENT_PASS
}

func (f *CFrame) handleGainedFocus(_ []interface{}, _ ...interface{}) cdk.EventFlag {
	// f.Lock()
	// defer f.Unlock()
	theme := f.GetTheme()
	if f.focusWithChild {
		if child := f.GetChild(); child != nil && child.IsFocus() {
			theme.Content.Normal = theme.Content.Focused
			theme.Border.Normal = theme.Border.Focused
		}
	} else {
		theme.Content.Normal = theme.Content.Focused
		theme.Border.Normal = theme.Border.Focused
	}
	f.SetThemeRequest(theme)
	f.Invalidate()
	return cdk.EVENT_PASS
}

// Text of the frame's label.
// Flags: Read / Write
// Default value: NULL
// const PropertyLabel cdk.Property = "label"

// A widget to display in place of the usual frame label.
// Flags: Read / Write
const PropertyLabelWidget cdk.Property = "label-widget"

// The horizontal alignment of the label.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0
const PropertyLabelXAlign cdk.Property = "label-x-align"

// The vertical alignment of the label.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
const PropertyLabelYAlign cdk.Property = "label-y-align"

// Deprecated property, use shadow_type instead.
// Flags: Read / Write
// Default value: GTK_SHADOW_ETCHED_IN
const PropertyShadow cdk.Property = "shadow"

// Appearance of the frame border.
// Flags: Read / Write
// Default value: GTK_SHADOW_ETCHED_IN
const PropertyShadowType cdk.Property = "shadow-type"
