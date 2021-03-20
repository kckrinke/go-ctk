// +build exampleHelloWorld

package main

import (
	"os"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-ctk"
)

func main() {
	// Construct a new CDK application
	app := cdk.NewApp(
		// program binary name
		"hello-world",
		// usage summary
		"the most basic CTK application",
		// because versioning is important
		"0.0.1",
		// used in logs, internal debugging, etc
		"helloWorld",
		// used where human-readable titles are necessary
		"Hello World",
		// the TTY device to use, /dev/tty is the default
		"/dev/tty",
		// initialize the user-interface
		func(d cdk.DisplayManager) error {
			// tell the display to listen for CTRL+C and interrupt gracefully
			d.CaptureCtrlC()
			// create a new window, give it a human-readable title
			w := ctk.NewWindowWithTitle("Hello World")
			// tell CTK that the window and it's contents are to be drawn upon
			// the terminal display
			w.ShowAll()
			// tell CDK that this window is the foreground window
			d.SetActiveWindow(w)
			// no errors to report, nil to proceed
			return nil
		},
	)
	// run the application, handing over the command-line arguments received
	if err := app.Run(os.Args); err != nil {
		// doesn't have to be a Fatal exit
		cdk.Fatal(err)
	}
	// end of program
}
