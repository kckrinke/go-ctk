package ctk

// TODO: new from stock id
// TODO: mnemonics support, GtkAccel?

import (
	"fmt"
	"strings"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
)

// CDK type-tag for Button objects
const TypeButton cdk.CTypeTag = "ctk-button"

var (
	DefaultMonoButtonTheme = cdk.Theme{
		Content: cdk.ThemeAspect{
			Normal:      cdk.DefaultMonoStyle,
			Focused:     cdk.DefaultMonoStyle.Dim(false),
			Active:      cdk.DefaultMonoStyle.Dim(false).Bold(true),
			FillRune:    cdk.DefaultFillRune,
			BorderRunes: cdk.DefaultBorderRune,
			ArrowRunes:  cdk.DefaultArrowRune,
			Overlay:     false,
		},
		Border: cdk.ThemeAspect{
			Normal:      cdk.DefaultMonoStyle,
			Focused:     cdk.DefaultMonoStyle.Dim(false),
			Active:      cdk.DefaultMonoStyle.Dim(false).Bold(true),
			FillRune:    cdk.DefaultFillRune,
			BorderRunes: cdk.DefaultBorderRune,
			ArrowRunes:  cdk.DefaultArrowRune,
			Overlay:     false,
		},
	}
	DefaultColorButtonTheme = cdk.Theme{
		Content: cdk.ThemeAspect{
			Normal:      cdk.DefaultColorStyle.Foreground(cdk.ColorWhite).Background(cdk.ColorFireBrick).Dim(true).Bold(false),
			Focused:     cdk.DefaultColorStyle.Foreground(cdk.ColorWhite).Background(cdk.ColorDarkRed).Dim(false).Bold(true),
			Active:      cdk.DefaultColorStyle.Foreground(cdk.ColorWhite).Background(cdk.ColorDarkRed).Dim(false).Bold(true).Reverse(true),
			FillRune:    cdk.DefaultFillRune,
			BorderRunes: cdk.DefaultBorderRune,
			ArrowRunes:  cdk.DefaultArrowRune,
			Overlay:     false,
		},
		Border: cdk.ThemeAspect{
			Normal:      cdk.DefaultColorStyle.Foreground(cdk.ColorWhite).Background(cdk.ColorFireBrick).Dim(true).Bold(false),
			Focused:     cdk.DefaultColorStyle.Foreground(cdk.ColorWhite).Background(cdk.ColorDarkRed).Dim(false).Bold(true),
			Active:      cdk.DefaultColorStyle.Foreground(cdk.ColorWhite).Background(cdk.ColorDarkRed).Dim(false).Bold(true).Reverse(true),
			FillRune:    cdk.DefaultFillRune,
			BorderRunes: cdk.DefaultBorderRune,
			ArrowRunes:  cdk.DefaultArrowRune,
			Overlay:     false,
		},
	}
)

func init() {
	_ = cdk.TypesManager.AddType(TypeButton, func() interface{} { return MakeButton() })
}

// Button Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- Button
//	          +- ToggleButton
//	          +- ColorButton
//	          +- FontButton
//	          +- LinkButton
//	          +- OptionMenu
//	          +- ScaleButton
type Button interface {
	Bin
	Activatable
	Alignable
	Buildable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	Activate() (value bool)
	Clicked() cdk.EventFlag
	SetRelief(newStyle ReliefStyle)
	GetRelief() (value ReliefStyle)
	GetLabel() (value string)
	SetLabel(label string)
	GetUseStock() (value bool)
	SetUseStock(useStock bool)
	GetUseUnderline() (value bool)
	SetUseUnderline(useUnderline bool)
	SetFocusOnClick(focusOnClick bool)
	GetFocusOnClick() (value bool)
	SetAlignment(xAlign float64, yAlign float64)
	GetAlignment() (xAlign float64, yAlign float64)
	SetImage(image Widget)
	GetImage() (value Widget, ok bool)
	SetImagePosition(position PositionType)
	GetImagePosition() (value PositionType)
	Add(w Widget)
	Remove(w Widget)
	SetPressed(pressed bool)
	GetPressed() bool
	GrabFocus()
	GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool)
	GetDefaultChildren() []Widget
	GetWidgetAt(p *cdk.Point2I) Widget
	CancelEvent()
	GrabEventFocus()
	ProcessEvent(evt cdk.Event) cdk.EventFlag
	Invalidate() cdk.EventFlag
	GetThemeRequest() (theme cdk.Theme)
	GetSizeRequest() (width, height int)
	Resize() cdk.EventFlag
	Draw(canvas cdk.Canvas) cdk.EventFlag
}

// The CButton structure implements the Button interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Button objects
type CButton struct {
	CBin

	pressed bool
	canvas  *cdk.CCanvas
}

// Default constructor for Button objects
func MakeButton() *CButton {
	b := NewButtonWithLabel("")
	return b
}

// Constructor for Button objects
func NewButton() *CButton {
	b := new(CButton)
	b.Init()
	return b
}

func NewButtonWithLabel(text string) (b *CButton) {
	b = NewButton()
	label := NewLabel(text)
	b.Add(label)
	label.SetTheme(DefaultColorButtonTheme)
	label.UnsetFlags(CAN_FOCUS)
	label.UnsetFlags(CAN_DEFAULT)
	label.UnsetFlags(RECEIVES_DEFAULT)
	label.SetLineWrap(false)
	label.SetLineWrapMode(cdk.WRAP_NONE)
	label.SetJustify(cdk.JUSTIFY_CENTER)
	label.SetAlignment(0.5, 0.5)
	label.SetSingleLineMode(true)
	label.Show()
	return b
}

// Creates a new Button containing a label. If characters in label are
// preceded by an underscore, they are underlined. If you need a literal
// underscore character in a label, use '__' (two underscores). The first
// underlined character represents a keyboard accelerator called a mnemonic.
// Pressing Alt and that key activates the button.
// Parameters:
// 	label	The text of the button, with an underscore in front of the
// mnemonic character
func NewButtonWithMnemonic(text string) (b *CButton) {
	b = NewButtonWithLabel(text)
	b.SetUseUnderline(true)
	return b
}

// Creates a new Button containing the image and text from a stock item.
// Some stock ids have preprocessor macros like GTK_STOCK_OK and
// GTK_STOCK_APPLY. If stock_id is unknown, then it will be treated as a
// mnemonic label (as for NewWithMnemonic).
// Parameters:
// 	stockId	the name of the stock item
// Returns:
// 	a new Button
func NewButtonFromStock(stockId StockID) (value *CButton) {
	b := NewButtonWithLabel(string(stockId))
	b.Init()
	if item := LookupStockItem(stockId); item != nil {
		b.SetUseStock(true)
		b.SetUseUnderline(true)
		b.SetLabel(item.Label)
	} else {
		b.SetLabel(string(stockId))
	}
	return b
}

// Constructor with Widget for Button objects, uncertain this works as expected
// due to struct type information loss on interface filter
func NewButtonWithWidget(w Widget) *CButton {
	b := new(CButton)
	b.Init()
	b.Add(w)
	return b
}

// Button object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Button instance
func (b *CButton) Init() (already bool) {
	if b.InitTypeItem(TypeButton, b) {
		return true
	}
	b.CBin.Init()
	b.flags = NULL_WIDGET_FLAG
	b.SetFlags(SENSITIVE | PARENT_SENSITIVE)
	b.SetFlags(CAN_DEFAULT | RECEIVES_DEFAULT | CAN_FOCUS)
	b.SetFlags(APP_PAINTABLE)
	b.SetTheme(DefaultColorButtonTheme)
	b.canvas = cdk.NewCanvas(cdk.MakePoint2I(0, 0), cdk.MakeRectangle(0, 0), b.GetTheme().Content.Normal)
	handle := fmt.Sprintf("%v.focus-changed", b.ObjectName())
	b.Connect(SignalLostFocus, handle, b.handleLostFocus)
	b.Connect(SignalGainedFocus, handle, b.handleGainedFocus)
	b.pressed = false
	_ = b.InstallBuildableProperty(PropertyFocusOnClick, cdk.BoolProperty, true, true)
	_ = b.InstallBuildableProperty(PropertyButtonLabel, cdk.StringProperty, true, nil)
	_ = b.InstallBuildableProperty(PropertyRelief, cdk.StructProperty, true, nil)
	_ = b.InstallBuildableProperty(PropertyUseStock, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyUseUnderline, cdk.BoolProperty, true, false)
	_ = b.InstallBuildableProperty(PropertyXAlign, cdk.FloatProperty, true, 0.5)
	_ = b.InstallBuildableProperty(PropertyYAlign, cdk.FloatProperty, true, 0.5)
	b.Connect(cdk.SignalSetProperty, fmt.Sprintf("%v.set-label", b.ObjectName()), func(data []interface{}, argv ...interface{}) cdk.EventFlag {
		if len(argv) == 3 {
			if key, ok := argv[1].(cdk.Property); ok {
				switch key {
				case PropertyButtonLabel:
					if val, ok := argv[2].(string); ok {
						b.SetLabel(val)
					} else {
						b.LogError("property label value is not string: %T", argv[2])
					}
				}
			}
		}
		// allow property to be set
		return cdk.EVENT_PASS
	})
	b.Invalidate()
	return false
}

func (b *CButton) Build(builder Builder, element *CBuilderElement) error {
	b.Freeze()
	defer b.Thaw()
	if name, ok := element.Attributes["id"]; ok {
		b.SetName(name)
	}
	if v, ok := element.Properties[PropertyUseStock.String()]; ok {
		b.SetUseStock(utils.IsTrue(v))
		b.SetUseUnderline(true)
	}
	if v, ok := element.Properties[PropertyLabel.String()]; ok {
		b.SetLabel(v)
	}
	for k, v := range element.Properties {
		switch cdk.Property(k) {
		case PropertyLabel:
		case PropertyUseStock:
		default:
			element.ApplyProperty(k, v)
		}
	}
	element.ApplySignals()
	return nil
}

func (b *CButton) Activate() (value bool) {
	return b.Emit(SignalActivate, b) == cdk.EVENT_STOP
}

// TODO: button Clicked() is not defined well

func (b *CButton) Clicked() cdk.EventFlag {
	return b.Emit(SignalClicked, b)
}

func (b *CButton) SetRelief(newStyle ReliefStyle) {
	if err := b.SetStructProperty(PropertyRelief, newStyle); err != nil {
		b.LogErr(err)
	}
}

func (b *CButton) GetRelief() (value ReliefStyle) {
	if v, err := b.GetStructProperty(PropertyRelief); err != nil {
		b.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(ReliefStyle); !ok {
			b.LogError("value stored in relief property is not of type ReliefStyle")
		}
	}
	return
}

// Fetches the text from the label of the button, as set by
// SetLabel. If the label text has not been set the return
// value will be NULL. This will be the case if you create an empty button
// with New to use as a container.
// Returns:
// 	The text of the label widget. This string is owned by the
// 	widget and must not be modified or freed.
func (b *CButton) GetLabel() (value string) {
	if v, ok := b.GetChild().(Label); ok {
		return v.GetText()
	}
	var err error
	if value, err = b.GetStringProperty(PropertyButtonLabel); err != nil {
		b.LogErr(err)
	}
	return
}

// Sets the text of the label of the button to str . This text is also used
// to select the stock item if SetUseStock is used. This will
// also clear any previously set labels.
// Parameters:
// 	label	a string
func (b *CButton) SetLabel(label string) {
	if b.GetUseStock() && label != "" {
		label = strings.ReplaceAll(label, "gtk", "ctk")
		if item := LookupStockItem(StockID(label)); item != nil {
			label = item.Label
		}
	}
	if v, ok := b.GetChild().(Label); ok {
		if strings.HasPrefix(label, "<markup") {
			if err := v.SetMarkup(label); err != nil {
				b.LogErr(err)
			}
		} else {
			v.SetText(label)
		}
	}
}

// Returns whether the button label is a stock item.
// Returns:
// 	TRUE if the button label is used to select a stock item instead
// 	of being used directly as the label text.
func (b *CButton) GetUseStock() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyUseStock); err != nil {
		b.LogErr(err)
	}
	return
}

// If TRUE, the label set on the button is used as a stock id to select the
// stock item for the button.
// Parameters:
// 	useStock	TRUE if the button should use a stock item
func (b *CButton) SetUseStock(useStock bool) {
	if err := b.SetBoolProperty(PropertyUseStock, useStock); err != nil {
		b.LogErr(err)
	} else {
		label := b.GetLabel()
		b.SetLabel(label)
	}
}

// Returns whether an embedded underline in the button label indicates a
// mnemonic. See SetUseUnderline.
// Returns:
// 	TRUE if an embedded underline in the button label indicates the
// 	mnemonic accelerator keys.
func (b *CButton) GetUseUnderline() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyUseUnderline); err != nil {
		b.LogErr(err)
	}
	return
}

// If true, an underline in the text of the button label indicates the next
// character should be used for the mnemonic accelerator key.
// Parameters:
// 	useUnderline	TRUE if underlines in the text indicate mnemonics
func (b *CButton) SetUseUnderline(useUnderline bool) {
	if err := b.SetBoolProperty(PropertyUseUnderline, useUnderline); err != nil {
		b.LogErr(err)
	}
	if child := b.GetChild(); child != nil {
		if label, ok := child.(Label); ok {
			label.SetUseUnderline(useUnderline)
		}
	}
}

// Sets whether the button will grab focus when it is clicked with the mouse.
// Making mouse clicks not grab focus is useful in places like toolbars where
// you don't want the keyboard focus removed from the main area of the
// application.
// Parameters:
// 	focusOnClick	whether the button grabs focus when clicked with the mouse
func (b *CButton) SetFocusOnClick(focusOnClick bool) {
	if err := b.SetBoolProperty(PropertyFocusOnClick, focusOnClick); err != nil {
		b.LogErr(err)
	}
}

// Returns whether the button grabs focus when it is clicked with the mouse.
// See SetFocusOnClick.
// Returns:
// 	TRUE if the button grabs focus when it is clicked with the
// 	mouse.
func (b *CButton) GetFocusOnClick() (value bool) {
	var err error
	if value, err = b.GetBoolProperty(PropertyFocusOnClick); err != nil {
		b.LogErr(err)
	}
	return
}

// Sets the alignment of the child. This property has no effect unless the
// child is a Misc or a Alignment.
// Parameters:
// 	xAlign	the horizontal position of the child, 0.0 is left aligned,
// 1.0 is right aligned
// 	yAlign	the vertical position of the child, 0.0 is top aligned,
// 1.0 is bottom aligned
func (b *CButton) SetAlignment(xAlign float64, yAlign float64) {
	xAlign = utils.ClampF(xAlign, 0.0, 1.0)
	yAlign = utils.ClampF(yAlign, 0.0, 1.0)
	if err := b.SetProperty(PropertyXAlign, xAlign); err != nil {
		b.LogErr(err)
	}
	if err := b.SetProperty(PropertyYAlign, yAlign); err != nil {
		b.LogErr(err)
	}
}

// Gets the alignment of the child in the button.
// Parameters:
// 	xAlign	return location for horizontal alignment.
// 	yAlign	return location for vertical alignment.
func (b *CButton) GetAlignment() (xAlign float64, yAlign float64) {
	var err error
	if xAlign, err = b.GetFloatProperty(PropertyXAlign); err != nil {
		b.LogErr(err)
	}
	err = nil
	if yAlign, err = b.GetFloatProperty(PropertyYAlign); err != nil {
		b.LogErr(err)
	}
	return
}

// Set the image of button to the given widget. Note that it depends on the
// gtk-button-images setting whether the image will be displayed or
// not, you don't have to call WidgetShow on image yourself.
// Parameters:
// 	image	a widget to set as the image for the button
func (b *CButton) SetImage(image Widget) {
	if err := b.SetStructProperty(PropertyImage, image); err != nil {
		b.LogErr(err)
	}
}

// Gets the widget that is currently set as the image of button . This may
// have been explicitly set by SetImage or constructed by
// NewFromStock.
// Returns:
// 	a Widget or NULL in case there is no image.
// 	[transfer none]
func (b *CButton) GetImage() (value Widget, ok bool) {
	if w, err := b.GetStructProperty(PropertyImage); err != nil {
		b.LogErr(err)
	} else {
		if value, ok = w.(Widget); !ok {
			value = nil
			return
		}
	}
	return
}

// Sets the position of the image relative to the text inside the button.
// Parameters:
// 	position	the position
func (b *CButton) SetImagePosition(position PositionType) {
	if err := b.SetStructProperty(PropertyImagePosition, position); err != nil {
		b.LogErr(err)
	}
}

// Gets the position of the image relative to the text inside the button.
// Returns:
// 	the position
func (b *CButton) GetImagePosition() (value PositionType) {
	if v, err := b.GetStructProperty(PropertyImagePosition); err != nil {
		b.LogErr(err)
	} else {
		var ok bool
		if value, ok = v.(PositionType); !ok {
			b.LogError("value stored in PropertyImagePosition is not a PositionType: %T", v)
		}
	}
	return
}

func (b *CButton) Add(w Widget) {
	if len(b.children) == 0 {
		b.CBin.Add(w)
		b.Invalidate()
	} else {
		b.LogError("button bin is full, failed to add: %v", w.ObjectName())
	}
}

func (b *CButton) Remove(w Widget) {
	if len(b.children) > 0 {
		b.CBin.Remove(w)
		b.Invalidate()
	} else {
		b.LogError("button bin is empty, failed to remove: %v", w.ObjectName())
	}
}

func (b *CButton) SetPressed(pressed bool) {
	b.pressed = pressed
	b.Invalidate()
	if pressed {
		b.Emit(SignalPressed)
	} else {
		b.Emit(SignalReleased)
	}
}

func (b *CButton) GetPressed() bool {
	return b.pressed
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
func (b *CButton) GrabFocus() {
	if b.CanFocus() {
		if r := b.Emit(SignalGrabFocus, b); r == cdk.EVENT_PASS {
			tl := b.GetWindow()
			if tl != nil {
				var fw Widget
				focused := tl.GetFocus()
				tl.SetFocus(b)
				if focused != nil {
					var ok bool
					if fw, ok = focused.(Widget); ok && fw.ObjectID() != b.ObjectID() {
						if f := fw.Emit(SignalLostFocus, fw); f == cdk.EVENT_STOP {
							fw = nil
						}
					}
				}
				if f := b.Emit(SignalGainedFocus, b, fw); f == cdk.EVENT_STOP {
					if fw != nil {
						tl.SetFocus(fw)
					}
				}
				b.LogDebug("has taken focus")
			}
		}
	}
}

func (b *CButton) GetFocusChain() (focusableWidgets []interface{}, explicitlySet bool) {
	focusableWidgets = []interface{}{b}
	return
}

func (b *CButton) GetDefaultChildren() []Widget {
	return []Widget{b}
}

func (b *CButton) GetWidgetAt(p *cdk.Point2I) Widget {
	if b.HasPoint(p) && b.IsVisible() {
		return b
	}
	return nil
}

func (b *CButton) CancelEvent() {
	if f := b.Emit(SignalCancelEvent, b); f == cdk.EVENT_PASS {
		b.SetPressed(false)
		b.ReleaseEventFocus()
	}
}

func (b *CButton) GrabEventFocus() {
	if window := b.GetWindow(); window != nil {
		if f := b.Emit(SignalGrabEventFocus, b, window); f == cdk.EVENT_PASS {
			window.SetEventFocus(b)
		}
	}
}

func (b *CButton) ProcessEvent(evt cdk.Event) cdk.EventFlag {
	switch e := evt.(type) {
	case *cdk.EventMouse:
		pos := cdk.NewPoint2I(e.Position())
		switch e.State() {
		case cdk.BUTTON_PRESS, cdk.DRAG_START:
			if b.HasPoint(pos) {
				if focusOnClick, err := b.GetBoolProperty(PropertyFocusOnClick); err == nil && focusOnClick {
					b.GrabFocus()
				}
				b.GrabEventFocus()
				b.SetPressed(true)
				b.LogDebug("pressed")
				return cdk.EVENT_STOP
			}
		case cdk.MOUSE_MOVE, cdk.DRAG_MOVE:
			if b.HasEventFocus() {
				if !b.HasPoint(pos) {
					b.LogDebug("out of bounds")
					b.CancelEvent()
					return cdk.EVENT_STOP
				}
			}
			return cdk.EVENT_PASS
		case cdk.BUTTON_RELEASE, cdk.DRAG_STOP:
			if b.HasEventFocus() {
				if !b.HasPoint(pos) {
					b.LogDebug("out of bounds")
					b.CancelEvent()
					return cdk.EVENT_STOP
				}
				b.ReleaseEventFocus()
				if f := b.Clicked(); f == cdk.EVENT_PASS {
					b.Activate()
				}
				b.SetPressed(false)
				b.LogDebug("released")
				return cdk.EVENT_STOP
			}
		}
	case *cdk.EventKey:
		if b.HasEventFocus() {
			b.LogDebug("keypress cancelling mouse event handling")
			b.CancelEvent()
			return cdk.EVENT_STOP
		}
		switch e.Key() {
		case cdk.KeyRune:
			if e.Rune() != ' ' {
				break
			}
			fallthrough
		case cdk.KeyEnter:
			if focusOnClick, err := b.GetBoolProperty(PropertyFocusOnClick); err == nil && focusOnClick {
				b.GrabFocus()
			}
			b.LogTrace("pressed")
			b.SetPressed(true)
			if f := b.Clicked(); f == cdk.EVENT_PASS {
				b.Activate()
			}
			b.SetPressed(false)
			b.LogTrace("released")
			return cdk.EVENT_STOP
		}
	}
	return cdk.EVENT_PASS
}

func (b *CButton) Invalidate() cdk.EventFlag {
	theme := b.GetThemeRequest()
	if child := b.GetChild(); child != nil {
		alloc := child.GetAllocation()
		local := child.GetOrigin()
		local.SubPoint(b.GetOrigin())
		if b.canvas == nil {
			b.canvas = cdk.NewCanvas(local, alloc, theme.Content.Normal)
		} else {
			b.canvas.SetOrigin(local)
			b.canvas.Resize(alloc, theme.Content.Normal)
		}
		child.SetTheme(theme)
		child.Invalidate()
		return cdk.EVENT_STOP
	}
	alloc := b.GetAllocation()
	if b.canvas == nil {
		b.canvas = cdk.NewCanvas(cdk.MakePoint2I(0, 0), alloc, theme.Content.Normal)
	} else {
		b.canvas.SetOrigin(cdk.MakePoint2I(0, 0))
		b.canvas.Resize(alloc, theme.Content.Normal)
	}
	return cdk.EVENT_STOP
}

func (b *CButton) GetThemeRequest() (theme cdk.Theme) {
	theme = b.CWidget.GetThemeRequest()
	if b.GetPressed() {
		theme.Content.Normal = theme.Content.Active
		theme.Content.Focused = theme.Content.Active
		theme.Border.Normal = theme.Border.Active
		theme.Border.Focused = theme.Border.Active
	} else if b.IsFocused() {
		theme.Content.Normal = theme.Content.Focused
		theme.Border.Normal = theme.Border.Focused
	}
	return
}

func (b *CButton) getBorderRequest() (border bool) {
	border = true
	alloc := b.GetAllocation()
	if alloc.W <= 2 || alloc.H <= 2 {
		border = false
	}
	return
}

func (b *CButton) GetSizeRequest() (width, height int) {
	size := cdk.NewRectangle(b.CWidget.GetSizeRequest())
	if child := b.GetChild(); child != nil {
		labelSizeReq := cdk.NewRectangle(child.GetSizeRequest())
		if size.W <= -1 && labelSizeReq.W > -1 {
			size.W = 2 + labelSizeReq.W + 2 // borders and bookends
		}
		if size.H <= -1 && labelSizeReq.H > -1 {
			size.H = labelSizeReq.H + 2 // borders
		}
	}
	return size.W, size.H
}

func (b *CButton) Resize() cdk.EventFlag {
	// our allocation has been set prior to Resize() being called
	child := b.GetChild()
	if child != nil {
		alloc := b.GetAllocation()
		if alloc.W <= 0 && alloc.H <= 0 {
			child.SetAllocation(cdk.MakeRectangle(0, 0))
			return child.Resize()
		}
		x, y := 0, 0
		origin := b.GetOrigin()
		if alloc.W >= 3 && alloc.H >= 3 {
			x, y = 1, 1
			alloc.Sub(2, 2)
		}
		childSize := cdk.NewRectangle(child.GetSizeRequest())
		childSize.W = alloc.W
		childSize.H = alloc.H
		childSize.Floor(0, 0)
		if label, ok := child.(Label); ok {
			maxChars, lineCount := label.GetPlainTextInfo()
			xAlign := 0.5
			yAlign := 0.5
			if lineCount >= childSize.H {
				yAlign = 0
			}
			if maxChars >= childSize.W {
				label.SetJustify(cdk.JUSTIFY_LEFT)
				xAlign = 0
			} else {
				label.SetJustify(cdk.JUSTIFY_CENTER)
			}
			label.SetAlignment(xAlign, yAlign)
		}
		local := cdk.MakePoint2I(x, y)
		theme := b.GetThemeRequest()
		if b.canvas == nil {
			b.canvas = cdk.NewCanvas(local, *childSize, theme.Content.Normal)
		} else {
			b.canvas.SetOrigin(local)
			b.canvas.Resize(*childSize, theme.Content.Normal)
		}
		child.SetTheme(theme)
		child.SetOrigin(origin.X+x, origin.Y+y)
		child.SetAllocation(*childSize)
		child.Resize()
	}
	b.Invalidate()
	return cdk.EVENT_PASS
}

func (b *CButton) Draw(canvas cdk.Canvas) cdk.EventFlag {
	b.Lock()
	defer b.Unlock()
	size := b.GetAllocation()
	if !b.IsVisible() || size.W <= 0 || size.H <= 0 {
		b.LogTrace("Draw(%v): not visible, zero width or zero height", canvas)
		canvas.Fill(b.GetTheme())
		return cdk.EVENT_STOP
	}

	var child Widget
	var label Label
	if child = b.GetChild(); child == nil {
		b.LogError("button child (label) not found")
		return cdk.EVENT_PASS
	} else if v, ok := child.(Label); ok {
		label = v
	}

	theme := b.GetThemeRequest()
	border := b.getBorderRequest()

	canvas.Box(
		cdk.MakePoint2I(0, 0),
		cdk.MakeRectangle(size.W, size.H),
		border, true,
		theme.Content.Overlay,
		theme.Content.FillRune,
		theme.Content.Normal,
		theme.Border.Normal,
		theme.Border.BorderRunes,
	)

	if label == nil {
		child.Draw(b.canvas)
		if err := canvas.Composite(b.canvas); err != nil {
			b.LogError("composite error: %v", err)
		}
	} else {
		label.Draw(b.canvas)
		if err := canvas.Composite(b.canvas); err != nil {
			b.LogError("composite error: %v", err)
		}
	}

	if debug, _ := b.GetBoolProperty(cdk.PropertyDebug); debug {
		canvas.DebugBox(cdk.ColorRed, b.ObjectInfo())
	}
	return cdk.EVENT_STOP
}

func (b *CButton) handleLostFocus(data []interface{}, argv ...interface{}) cdk.EventFlag {
	_ = b.Invalidate()
	return cdk.EVENT_PASS
}

func (b *CButton) handleGainedFocus(data []interface{}, argv ...interface{}) cdk.EventFlag {
	_ = b.Invalidate()
	return cdk.EVENT_PASS
}

// Whether the button grabs focus when it is clicked with the mouse.
// Flags: Read / Write
// Default value: TRUE
const PropertyFocusOnClick cdk.Property = "focus-on-click"

// Child widget to appear next to the button text.
// Flags: Read / Write
const PropertyImage cdk.Property = "image"

// The position of the image relative to the text inside the button.
// Flags: Read / Write
// Default value: GTK_POS_LEFT
const PropertyImagePosition cdk.Property = "image-position"

// Text of the label widget inside the button, if the button contains a label
// widget.
// Flags: Read / Write / Construct
// Default value: NULL
const PropertyButtonLabel cdk.Property = "label"

// The border relief style.
// Flags: Read / Write
// Default value: GTK_RELIEF_NORMAL
const PropertyRelief cdk.Property = "relief"

// If set, the label is used to pick a stock item instead of being displayed.
// Flags: Read / Write / Construct
// Default value: FALSE
const PropertyUseStock cdk.Property = "use-stock"

// If set, an underline in the text indicates the next character should be
// used for the mnemonic accelerator key.
// Flags: Read / Write / Construct
// Default value: FALSE
const PropertyUseUnderline cdk.Property = "use-underline"

// If the child of the button is a Misc or Alignment, this property can
// be used to control it's horizontal alignment. 0.0 is left aligned, 1.0 is
// right aligned.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
// const PropertyXAlign cdk.Property = "xalign"

// If the child of the button is a Misc or Alignment, this property can
// be used to control it's vertical alignment. 0.0 is top aligned, 1.0 is
// bottom aligned.
// Flags: Read / Write
// Allowed values: [0,1]
// Default value: 0.5
// const PropertyYAlign cdk.Property = "yalign"

// The activate signal on Button is an action signal and emitting it causes
// the button to animate press then release. Applications should never
// connect to this signal, but use the clicked signal.
// const SignalActivate cdk.Signal = "activate"

// Emitted when the button has been activated (pressed and released).
const SignalClicked cdk.Signal = "clicked"

// Emitted when the pointer enters the button.
const SignalEnter cdk.Signal = "enter"

// Emitted when the pointer leaves the button.
const SignalLeave cdk.Signal = "leave"

// Emitted when the button is pressed.
const SignalPressed cdk.Signal = "pressed"

// Emitted when the button is released.
const SignalReleased cdk.Signal = "released"
