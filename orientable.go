package ctk

import (
	"github.com/kckrinke/go-cdk"
)

type Orientable interface {
	GetOrientation() (orientation cdk.Orientation)
	SetOrientation(orientation cdk.Orientation)
}

// The orientation of the orientable.
// Flags: Read / Write
// Default value: ORIENTATION_HORIZONTAL
const PropertyOrientation cdk.Property = "orientation"
