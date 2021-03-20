package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gobuffalo/envy"
	"golang.org/x/text/unicode/runenames"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
	"github.com/kckrinke/go-ctk"
)

const (
	AppName    = "ctk-app"
	AppUsage   = "an example of a CLI application using CTK"
	AppVersion = "0.0.1"
	AppTag     = "demo"
	AppTitle   = "CTK Demo"
)

// Build Configuration Flags
// use `go build -v -ldflags="-X 'main.IncludeLogFullPaths=false'"
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
	app := cdk.NewApp(
		AppName,
		AppUsage,
		AppVersion,
		AppTag,
		AppTitle,
		"/dev/tty",
		setup,
	)
	if err := app.Run(os.Args); err != nil {
		cdk.Fatal(err)
	}
}

func setup(d cdk.DisplayManager) error {
	theme := cdk.DefaultColorTheme
	theme.Content.Normal = theme.Content.Focused.Dim(false)
	theme.Border.Normal = theme.Border.Focused.Dim(false)
	d.CaptureCtrlC()
	w := ctk.NewWindowWithTitle("Character Map")
	w.SetTheme(theme)
	w.Connect(ctk.SignalEventKey, "escape-quit", func(data []interface{}, argv ...interface{}) cdk.EventFlag {
		if evt, ok := argv[1].(cdk.Event); ok {
			switch e := evt.(type) {
			case *cdk.EventKey:
				if e.Key() == cdk.KeyEscape {
					w.LogInfo("window caught escape key, quitting now")
					d.RequestQuit()
				}
			}
			return cdk.EVENT_STOP
		}
		return cdk.EVENT_PASS
	})
	align := ctk.NewAlignment(0.5, 0.5, 0.0, 0.0)
	align.SetTheme(theme)
	frame := ctk.NewFrame("Character Details")
	frame.SetTheme(theme)
	label := ctk.NewLabel("")
	label.SetJustify(cdk.JUSTIFY_LEFT)
	label.SetAlignment(0.5, 0.5)
	label.SetTheme(theme)
	label.SetText("loading...")
	label.SetBoolProperty("debug", false)
	label.SetSizeRequest(35, 4)
	label.Show()
	frame.Add(label)
	frame.Show()
	align.Add(frame)
	align.Show()
	w.GetVBox().PackStart(align, true, true, 0)
	ctx := d.App().GetContext()
	args := ctx.Args().Slice()
	cdk.AddTimeout(
		time.Millisecond*150,
		func() cdk.EventFlag {
			var message string
			if len(args) > 0 {
				if num, err := strconv.Atoi(args[0]); err != nil {
					message = fmt.Sprintf("invalid argument: %v", args[1])
				} else {
					r, w := utf8.DecodeRune([]byte(string(rune(num))))
					name := runenames.Name(r)
					if unicode.IsGraphic(r) {
						message = fmt.Sprintf("%s\nEntity: &#%d;\nPrint: \"%c\"\nWidth: \"%d\"", name, num, r, w)
					} else {
						message = fmt.Sprintf("%s\nEntity: &#%d;\nPrint: \"%x\"\nWidth: \"%d\"", name, num, r, w)
					}
				}
			} else {
				message = fmt.Sprintf(
					"Character Set:\t%v\nDisplay Size:\t%v\nDisplay Colors:\t%d\nTerminal:\t%v",
					d.Display().CharacterSet(),
					w.GetAllocation(),
					d.Display().Colors(),
					envy.Get("TERM", "(unset)"),
				)
			}
			label.SetText(message)
			d.RequestDraw()
			d.RequestShow()
			return cdk.EVENT_STOP
		},
	)
	w.ShowAll()
	d.SetActiveWindow(w)
	return nil
}
