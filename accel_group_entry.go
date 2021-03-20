package ctk

import (
	"github.com/kckrinke/go-cdk"
)

type AccelGroupEntry struct {
	Accelerator AccelKey
	Closure     GClosure
	Quark       cdk.QuarkID
}

func NewAccelGroupEntry(key AccelKey, closure GClosure, quark cdk.QuarkID) (age *AccelGroupEntry) {
	age = &AccelGroupEntry{
		Accelerator: key,
		Closure:     closure,
		Quark:       quark,
	}
	return
}
