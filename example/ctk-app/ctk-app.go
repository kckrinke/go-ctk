// +build exampleCtkApp

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
	"github.com/kckrinke/go-ctk"
)

const (
	APP_NAME    = "ctk-app"
	APP_USAGE   = "an example of a CLI application using CTK"
	APP_VERSION = "0.0.1"
	APP_TAG     = "demo"
	APP_TITLE   = "CTK Demo"
)

// Build Configuration Flags
// use `go build -v -ldflags="-X 'main.IncludeLogFullPaths=true'"
var (
	IncludeLogFullPaths  string = "true"
	IncludeLogTimestamps string = "false"
	IncludeProfiling     string = "false"
	Debug                bool   = false
)

func main() {
	cdk.Build.LogFullPaths = utils.IsTrue(IncludeLogFullPaths)
	cdk.Build.LogTimestamps = utils.IsTrue(IncludeLogTimestamps)
	cdk.Build.Profiling = utils.IsTrue(IncludeProfiling)
	app := cdk.NewApp(APP_NAME, APP_USAGE, APP_VERSION, APP_TAG, APP_TITLE, "/dev/tty", setupUi)
	app.AddFlag(&cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
	})
	app.AddFlag(&cli.StringFlag{
		Name:  "example-flag",
		Value: "testing",
	})
	app.AddCommand(&cli.Command{
		Name:  "demo-cmd",
		Usage: "demonstrate custom commands",
		Action: func(c *cli.Context) error {
			cdk.InfoF("demo-cmd command action")
			return nil
		},
	})
	if err := app.Run(os.Args); err != nil {
		cdk.Fatal(err)
	}
}

func setupUi(manager cdk.DisplayManager) error {
	if manager.App().GetContext().Bool("debug") {
		cdk.DebugF("enabling debug")
		Debug = true
	}
	// note that screen is captured at this time!
	manager.CaptureCtrlC()
	w := ctk.NewWindowWithTitle(APP_TITLE)
	w.Show()
	w.SetSensitive(true)
	manager.SetActiveWindow(w)
	vbox := w.GetVBox()
	vbox.SetHomogeneous(true)
	// vbox.SetBoolProperty("debug", true)
	b := newButton("b1", "Quit Button (expand,fill)", func(data []interface{}, argv ...interface{}) cdk.EventFlag {
		cdk.InfoF("Exiting now.")
		manager.RequestQuit()
		return cdk.EVENT_STOP
	})
	b.Show()
	vbox.PackStart(b, true, true, 0)

	// another row
	hbox2 := ctk.NewHBox(false, 0)
	hbox2.Show()
	if Debug {
		hbox2.SetBoolProperty("debug", true)
		hbox2.SetBoolProperty("debug-children", true)
	}
	vbox.PackStart(hbox2, true, true, 0)

	frame := ctk.NewFrame("This is a frame")
	frame.SetSizeRequest(30, -1)
	frame.SetFocusWithChild(true)
	frame.Show()
	hbox2.PackStart(frame, false, false, 0)
	// frame.SetLabelAlign(0.0, 0.5)
	if Debug {
		frame.SetBoolProperty("debug", true)
	}
	l1 := newLabel(IPSUM_LONG_MARKUP)
	l1.SetSizeRequest(35, -1)
	// l1.SetMaxWidthChars(35)
	l1.SetLineWrapMode(cdk.WRAP_CHAR)
	l1.SetJustify(cdk.JUSTIFY_LEFT)
	l1.SetSingleLineMode(false)
	if Debug {
		l1.SetBoolProperty("debug", true)
	}
	l1.Show()

	sv := ctk.NewScrolledViewport()
	sv.SetPolicy(ctk.PolicyAutomatic, ctk.PolicyAutomatic)
	sv.Show()
	sv.Add(l1)
	frame.Add(sv)

	hbox3 := ctk.NewHBox(false, 0)
	hbox3.Show()
	if Debug {
		hbox3.SetBoolProperty("debug", true)
		hbox3.SetBoolProperty("debug-children", true)
	}
	// hbox2.SetBoolProperty("debug-children", true)
	hbox2.PackStart(hbox3, true, true, 0)

	b2 := newButton("b2", "B2 (expand+fill)", func(d []interface{}, argv ...interface{}) cdk.EventFlag {
		cdk.InfoF("pressed button #2")
		return cdk.EVENT_STOP
	})
	b2.Show()
	hbox3.PackStart(b2, true, true, 0)

	b4 := newButton("curses", "Curses<u><i>!</i></u>", func(d []interface{}, argv ...interface{}) cdk.EventFlag {
		cdk.InfoF("pressed Curses!")
		dialog := ctk.NewDialogWithButtons("dialog title", w, ctk.DialogModal,
			ctk.StockOk, ctk.ResponseOk,
			ctk.StockCancel, ctk.ResponseCancel,
		)
		dialog.SetSizeRequest(40, 10)
		label := ctk.NewLabel("testing the content area")
		label.Show()
		dialog.GetContentArea().PackStart(label, true, true, 0)
		dialog.Show()
		// if Debug {
		// dialog.GetVBox().SetBoolProperty(cdk.PropertyDebug, true)
		// dialog.GetVBox().SetBoolProperty(Property, true)
		// }
		response := dialog.Run()
		go func() {
			select {
			case r := <-response:
				dialog.Destroy()
				_ = dialog.DestroyObject()
				cdk.DebugF("dialog response: %v", r)
			}
		}()
		return cdk.EVENT_STOP
	})
	b4.SetSizeRequest(13, 3)
	b4.Show()
	hbox3.PackEnd(b4, false, false, 0)
	b4.GrabFocus()

	b3 := newButton("b3", "B3 (expand)", func(d []interface{}, argv ...interface{}) cdk.EventFlag {
		cdk.InfoF("pressed button #3")
		return cdk.EVENT_STOP
	})
	// b3.SetSizeRequest(10, 3)
	b3.Show()
	hbox3.PackStart(b3, true, false, 0)
	return nil
}

var (
	// IPSUM_LONG_PLAIN = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt orci a quam dignissim mattis. Nulla volutpat egestas nibh vitae facilisis. Nam dictum risus a nisl suscipit, in luctus felis facilisis. Sed et ante pellentesque, vehicula dui vel, dictum eros. Duis convallis sem vitae tellus feugiat rhoncus. Curabitur risus lectus, elementum id molestie vel, gravida fermentum libero. In aliquet massa eu tellus pulvinar, in scelerisque ipsum ultricies. Quisque elementum nulla vitae condimentum venenatis. Vestibulum vitae lectus sit amet ipsum congue semper ornare tempus magna. Aliquam varius, eros eget ultrices auctor, lacus nibh blandit purus, sed rhoncus erat ex sed enim."
	IPSUM_LONG_MARKUP = "Lorem <i>ipsum</i> dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt orci a quam dignissim mattis. Nulla volutpat egestas nibh vitae facilisis. Nam dictum risus a nisl suscipit, in luctus felis facilisis. Sed et ante pellentesque, vehicula dui vel, dictum eros. Duis convallis sem vitae tellus feugiat rhoncus. Curabitur risus lectus, elementum id molestie vel, gravida fermentum libero. In aliquet massa eu tellus pulvinar, in scelerisque <i>ipsum</i> ultricies. Quisque elementum nulla vitae condimentum venenatis. Vestibulum vitae lectus sit amet <i>ipsum</i> congue semper ornare tempus magna. Aliquam varius, eros eget ultrices auctor, lacus nibh blandit purus, sed rhoncus erat ex sed enim."
)

func newLabel(text string) ctk.Label {
	l, err := ctk.NewLabelWithMarkup(text)
	if err != nil {
		cdk.Fatal(err)
	}
	if Debug {
		l.SetBoolProperty("debug", true)
	}
	return l
}

func newButton(name string, label string, fn cdk.SignalListenerFn) ctk.Button {
	b := ctk.NewButtonWithLabel("")
	if child := b.GetChild(); child != nil {
		if l, ok := child.(ctk.Label); ok {
			l.SetMarkup(label)
			l.SetEllipsize(true)
		}
	}
	b.SetName(name)
	b.SetSensitive(true)
	if Debug {
		b.SetBoolProperty("debug", true)
	}
	b.Connect(
		ctk.SignalActivate,
		fmt.Sprintf("%s.activate", name),
		fn,
	)
	return b
}

func newArrow(name string, arrow ctk.ArrowType, fn cdk.SignalListenerFn) ctk.Button {
	a := ctk.NewArrow(arrow)
	b := ctk.NewButtonWithWidget(a)
	b.SetSensitive(true)
	b.SetSizeRequest(1, 1)
	if Debug {
		b.SetBoolProperty("debug", true)
	}
	b.Connect(
		ctk.SignalActivate,
		fmt.Sprintf("%v.activate", name),
		fn,
	)
	return b
}
