package ctk

import (
	"github.com/kckrinke/go-cdk"
)

// CDK type-tag for EventBox objects
const TypeEventBox cdk.CTypeTag = "ctk-event-box"

func init() {
	_ = cdk.TypesManager.AddType(TypeEventBox, func() interface{} { return MakeEventBox() })
}

// EventBox Hierarchy:
//	Object
//	  +- Widget
//	    +- Container
//	      +- Bin
//	        +- EventBox
type EventBox interface {
	Bin
	Buildable

	Init() (already bool)
	SetAboveChild(aboveChild bool)
	GetAboveChild() (value bool)
	SetVisibleWindow(visibleWindow bool)
	GetVisibleWindow() (value bool)
	GrabFocus()
	Activate() (value bool)
	CancelEvent()
	ProcessEvent(evt cdk.Event) cdk.EventFlag
}

// The CEventBox structure implements the EventBox interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with EventBox objects
type CEventBox struct {
	CBin
}

// Default constructor for EventBox objects
func MakeEventBox() *CEventBox {
	return NewEventBox()
}

// Constructor for EventBox objects
func NewEventBox() (value *CEventBox) {
	e := new(CEventBox)
	e.Init()
	return e
}

// EventBox object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the EventBox instance
func (b *CEventBox) Init() (already bool) {
	if b.InitTypeItem(TypeEventBox, b) {
		return true
	}
	b.CBin.Init()
	b.SetFlags(SENSITIVE | PARENT_SENSITIVE)
	b.SetFlags(CAN_DEFAULT | RECEIVES_DEFAULT | CAN_FOCUS)
	b.SetFlags(APP_PAINTABLE)
	_ = b.InstallProperty(PropertyAboveChild, cdk.BoolProperty, true, false)
	_ = b.InstallProperty(PropertyVisibleWindow, cdk.BoolProperty, true, true)
	return false
}

// Set whether the event box window is positioned above the windows of its
// child, as opposed to below it. If the window is above, all events inside
// the event box will go to the event box. If the window is below, events in
// windows of child widgets will first got to that widget, and then to its
// parents. The default is to keep the window below the child.
// Parameters:
// 	aboveChild	TRUE if the event box window is above the windows of its child
func (e *CEventBox) SetAboveChild(aboveChild bool) {
	if err := e.SetBoolProperty(PropertyAboveChild, aboveChild); err != nil {
		e.LogErr(err)
	}
}

// Returns whether the event box window is above or below the windows of its
// child. See SetAboveChild for details.
// Returns:
// 	TRUE if the event box window is above the window of its child.
func (e *CEventBox) GetAboveChild() (value bool) {
	var err error
	if value, err = e.GetBoolProperty(PropertyAboveChild); err != nil {
		e.LogErr(err)
	}
	return
}

// Set whether the event box uses a visible or invisible child window. The
// default is to use visible windows. In an invisible window event box, the
// window that the event box creates is a GDK_INPUT_ONLY window, which means
// that it is invisible and only serves to receive events. A visible window
// event box creates a visible (GDK_INPUT_OUTPUT) window that acts as the
// parent window for all the widgets contained in the event box. You should
// generally make your event box invisible if you just want to trap events.
// Creating a visible window may cause artifacts that are visible to the
// user, especially if the user is using a theme with gradients or pixmaps.
// The main reason to create a non input-only event box is if you want to set
// the background to a different color or draw on it.
// Parameters:
// 	visibleWindow	boolean value
func (e *CEventBox) SetVisibleWindow(visibleWindow bool) {
	if err := e.SetBoolProperty(PropertyVisibleWindow, visibleWindow); err != nil {
		e.LogErr(err)
	}
}

// Returns whether the event box has a visible window. See
// SetVisibleWindow for details.
// Returns:
// 	TRUE if the event box window is visible
func (e *CEventBox) GetVisibleWindow() (value bool) {
	var err error
	if value, err = e.GetBoolProperty(PropertyVisibleWindow); err != nil {
		e.LogErr(err)
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
func (b *CEventBox) GrabFocus() {
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

func (b *CEventBox) Activate() (value bool) {
	return b.Emit(SignalActivate, b) == cdk.EVENT_STOP
}

func (b *CEventBox) CancelEvent() {
	b.Emit(SignalCancelEvent, b)
}

func (b *CEventBox) ProcessEvent(evt cdk.Event) cdk.EventFlag {
	return b.Emit(SignalCdkEvent, b, evt)
}

// Whether the event-trapping window of the eventbox is above the window of
// the child widget as opposed to below it.
// Flags: Read / Write
// Default value: FALSE
const PropertyAboveChild cdk.Property = "above-child"

// Whether the event box is visible, as opposed to invisible and only used to
// trap events.
// Flags: Read / Write
// Default value: TRUE
const PropertyVisibleWindow cdk.Property = "visible-window"
