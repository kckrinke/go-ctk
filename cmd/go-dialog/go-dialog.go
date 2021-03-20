package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
	"github.com/kckrinke/go-ctk"
)

const (
	AppName    = "go-dialog"
	AppUsage   = "display dialog boxes from shell scripts"
	AppVersion = "0.0.1"
	AppTag     = "go-dialog"
	AppTitle   = "go-dialog"
)

//go:embed dialog-msgbox.glade
var gladeMsgBox string

//go:embed dialog-yesno.glade
var gladeYesNo string

// Build Configuration Flags
//    go build -v -ldflags="-X 'main.IncludeLogTimestamps=true'"
var (
	IncludeLogFullPaths  = "true"
	IncludeLogTimestamps = "false"
	IncludeProfiling     = "false"
)

func init() {
	cdk.Build.LogFullPaths = utils.IsTrue(IncludeLogFullPaths)
	cdk.Build.LogTimestamps = utils.IsTrue(IncludeLogTimestamps)
	cdk.Build.Profiling = utils.IsTrue(IncludeProfiling)
}

func main() {
	app := cdk.NewApp(
		AppName, AppUsage, AppVersion,
		AppTag, AppTitle,
		"/dev/tty",
		setupUserInterface,
	)
	app.AddFlag(&cli.StringFlag{
		Name:  "back-title",
		Usage: "specify the window title text",
		Value: "",
	})
	app.AddFlag(&cli.StringFlag{
		Name:  "title",
		Usage: "specify the dialog title text",
		Value: "",
	})
	app.AddFlag(&cli.BoolFlag{
		Name:  "print-maxsize",
		Usage: "print the width and height on stdout and exit",
		Value: false,
	})
	app.AddCommand(&cli.Command{
		Name:      "msgbox",
		Usage:     "display a message with an OK button, each string following msgbox is a new line and concatenated into the message",
		ArgsUsage: "[message lines]",
		Action:    app.MainActionFn,
	})
	app.AddCommand(&cli.Command{
		Name:      "yesno",
		Usage:     "display a yes/no prompt with a message (see msgbox)",
		ArgsUsage: "[message lines]",
		Action:    app.MainActionFn,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "default",
				Usage:       "specify which button is focused initially",
				Value:       "yes",
				DefaultText: "yes",
			},
		},
	})
	if err := app.Run(os.Args); err != nil {
		cdk.Fatal(err)
	}
}

func setupUserInterface(dm cdk.DisplayManager) error {
	ctx := dm.App().GetContext()
	if ctx.Bool("print-maxsize") {
		if display := dm.Display(); display != nil {
			w, h := display.Size()
			dm.AddQuitHandler("print-maxsize", func() {
				fmt.Printf("%v %v\n", w, h)
			})
			dm.RequestQuit()
			return nil
		}
	}
	dm.LogInfo("setting up user interface")
	dm.CaptureCtrlC()

	builder := ctk.NewBuilder()
	var proceed bool
	switch ctx.Command.Name {
	case "msgbox":
		if err := setupUiMsgbox(ctx, builder, dm); err != nil {
			return err
		}
		proceed = true
	case "yesno":
		if err := setupUiYesNo(ctx, builder, dm); err != nil {
			return err
		}
		proceed = true
	case "":
		dm.AddQuitHandler("see-help", func() {
			fmt.Printf("see: %v --help\n", dm.App().Name())
		})
		dm.RequestQuit()
		return nil
	default:
		return fmt.Errorf("invalid command: %v", ctx.Command.Name)
	}
	if proceed {
		if err := startupUiDialog(ctx, builder, dm); err != nil {
			return err
		}
		dm.LogInfo("user interface set up complete")
		return nil
	}
	return fmt.Errorf("error intializing user interface")
}

func startupUiDialog(ctx *cli.Context, builder ctk.Builder, dm cdk.DisplayManager) error {
	backTitle := ctx.String("back-title")
	title := ctx.String("title")
	window := getWindow(builder)
	dialog := getDialog(builder)
	if window != nil {
		window.Show()
		window.SetTitle(backTitle)
		dm.SetActiveWindow(window)
		if dialog != nil {
			if display := dm.Display(); display != nil {
				dw, dh := display.Size()
				if dw > 22 && dh > 12 {
					sr := cdk.NewRectangle(dw/3, dh/3)
					sr.Clamp(20, 10, dw, dh)
					dialog.SetSizeRequest(sr.W, sr.H)
				}
				dialog.SetTransientFor(window)
				dialog.SetTitle(title)
			}
			dialog.Show()
			dialog.LogInfo("starting Run()")
			defBtn := ctx.String("default")
			switch strings.ToLower(defBtn) {
			case "yes", "ctk-yes":
				if yes := builder.GetWidget("yesno-yes"); yes != nil {
					if yw, ok := yes.(ctk.Widget); ok {
						yw.GrabFocus()
					}
				}
			case "no", "ctk-no":
				if no := builder.GetWidget("yesno-no"); no != nil {
					if nw, ok := no.(ctk.Widget); ok {
						nw.GrabFocus()
					}
				}
			}
			response := dialog.Run()
			go func() {
				select {
				case r := <-response:
					dialog.Destroy()
					_ = dialog.DestroyObject()
					dm.AddQuitHandler("dialog-response", func() {
						fmt.Printf("%v\n", r)
					})
					dm.RequestQuit()
				}
			}()
		} else {
			builder.LogError("missing main-dialog")
		}
	} else {
		builder.LogError("missing main-window")
	}
	return nil
}

func setupUiMsgbox(ctx *cli.Context, builder ctk.Builder, dm cdk.DisplayManager) error {
	builder.AddNamedSignalHandler("msgbox-ok", func(data []interface{}, argv ...interface{}) cdk.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(ctk.ResponseOk)
		} else {
			builder.LogError("msgbox-ok missing main-dialog")
		}
		return cdk.EVENT_STOP
	})
	if tmpl, err := template.New("msgbox").Parse(gladeMsgBox); err != nil || tmpl == nil {
		dm.LogErr(err)
	} else {
		content := ""
		if ctx.Args().Len() < 1 {
			return fmt.Errorf("msgbox missing message to display")
		}
		for i := 0; i < ctx.Args().Len(); i++ {
			if content != "" {
				content += "\n"
			}
			content += fmt.Sprintf("%v", ctx.Args().Get(i))
		}

		buff := new(bytes.Buffer)
		data := struct {
			Message string
		}{
			Message: content,
		}
		if err := tmpl.Execute(buff, data); err == nil {
			xml := string(buff.Bytes())
			var err error
			if _, err = builder.LoadFromString(xml); err != nil {
				return err
			}
		}
	}
	return nil
}

func setupUiYesNo(ctx *cli.Context, builder ctk.Builder, dm cdk.DisplayManager) error {
	builder.AddNamedSignalHandler("yesno-yes", func(data []interface{}, argv ...interface{}) cdk.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(ctk.ResponseYes)
		} else {
			builder.LogError("yesno-yes missing main-dialog")
		}
		return cdk.EVENT_STOP
	})
	builder.AddNamedSignalHandler("yesno-no", func(data []interface{}, argv ...interface{}) cdk.EventFlag {
		if dialog := getDialog(builder); dialog != nil {
			dialog.Response(ctk.ResponseNo)
		} else {
			builder.LogError("yesno-no missing main-dialog")
		}
		return cdk.EVENT_STOP
	})
	if tmpl, err := template.New("yesno").Parse(gladeYesNo); err != nil || tmpl == nil {
		dm.LogErr(err)
	} else {
		content := ""
		if ctx.Args().Len() < 1 {
			return fmt.Errorf("yesno missing message to display")
		}
		for i := 0; i < ctx.Args().Len(); i++ {
			if content != "" {
				content += "\n"
			}
			content += fmt.Sprintf("%v", ctx.Args().Get(i))
		}
		buff := new(bytes.Buffer)
		data := struct {
			Message string
		}{
			Message: content,
		}
		if err := tmpl.Execute(buff, data); err == nil {
			xml := string(buff.Bytes())
			var err error
			if _, err = builder.LoadFromString(xml); err != nil {
				return err
			}
		}
	}
	return nil
}

func getWindow(builder ctk.Builder) (window ctk.Window) {
	if mw := builder.GetWidget("main-window"); mw != nil {
		var ok bool
		if window, ok = mw.(ctk.Window); !ok {
			builder.LogError("main-window widget is not of ctk.Window type: %v (%T)", mw, mw)
		}
	} else {
		builder.LogError("missing main-window widget")
	}
	return
}

func getDialog(builder ctk.Builder) (dialog ctk.Dialog) {
	if md := builder.GetWidget("main-dialog"); md != nil {
		var ok bool
		if dialog, ok = md.(ctk.Dialog); !ok {
			builder.LogError("main-dialog widget is not of ctk.Dialog type: %v (%T)", md, md)
		}
	} else {
		builder.LogError("missing main-dialog widget")
	}
	return
}
