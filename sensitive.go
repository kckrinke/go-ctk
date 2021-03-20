package ctk

import (
	"github.com/kckrinke/go-cdk"
)

type Sensitive interface {
	Object

	GetWindow() Window
	CanFocus() bool
	IsFocus() bool
	IsFocused() bool
	IsVisible() bool
	GrabFocus()
	CancelEvent()
	IsSensitive() bool
	SetSensitive(sensitive bool)
	ProcessEvent(evt cdk.Event) cdk.EventFlag
}
