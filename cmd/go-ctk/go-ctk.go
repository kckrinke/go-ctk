package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
	"github.com/kckrinke/go-ctk/cmd/go-ctk/ctkfmt"
	"github.com/kckrinke/go-ctk/cmd/go-ctk/glade"
	"github.com/kckrinke/go-ctk/cmd/go-ctk/gtkdoc2ctk"
)

// Build Configuration Flags
// use `go build -v -ldflags="-X 'main.IncludeLogFullPaths=true'"
var (
	IncludeLogFullPaths  = "true"
	IncludeLogTimestamps = "false"
	IncludeProfiling     = "false"
)

func main() {
	cdk.Build.LogFullPaths = utils.IsTrue(IncludeLogFullPaths)
	cdk.Build.LogTimestamps = utils.IsTrue(IncludeLogTimestamps)
	cdk.Build.Profiling = utils.IsTrue(IncludeProfiling)
	app := &cli.App{
		Name:   "go-ctk",
		Usage:  "Curses Tool Kit Utility",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-file",
				Aliases:     []string{"log"},
				Usage:       "enable logging and specify log path to write to",
				Value:       os.TempDir() + string(os.PathSeparator) + "go-ctk.log",
				DefaultText: os.TempDir() + string(os.PathSeparator) + "go-ctk.log",
			},
		},
		Commands: []*cli.Command{
			ctkfmt.CliCommand,
			gtkdoc2ctk.CliCommand,
			glade.CliCommand,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		cdk.Fatal(err)
	}
}

func action(c *cli.Context) error {
	fmt.Println("this should display useful details about the current terminal environment")
	return nil
}
