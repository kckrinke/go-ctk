package ctk

import (
	"github.com/kckrinke/go-cdk"
)

func TestingWithCtkWindow(d cdk.DisplayManager) error {
	w := NewWindowWithTitle(d.GetTitle())
	d.SetActiveWindow(w)
	return nil
}
