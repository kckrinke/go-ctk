# CTK - Curses Tool Kit

Golang package to provide an advanced terminal user interface with a [GTK2]
inspired API, built upon [CDK].

## Notice

The current status of this project is in the v0.0.x versioning range and is
entirely experimental, or rather very unfinished and implemented poorly. There
were a number of fundamental misconceptions about how to approach a GTK-flavoured
curses programming experience. Please do not use this project for anything other
then intellectual curiosity. There's a lot here and not much is documented very
well and there's no "overview" of the codebase provided at this time.

This project is in the early process of a fundamental rewrite, in conjunction
with the [CDK] rewrite underway.

## Getting Started

CTK is a Go library and as such can be used in any of the typical Golang ways.

### Prerequisites

Go v1.16 (or later) is required in order to build and use the package. Beyond
that, there aren't any other dependencies. Visit: https://golang.org/doc/install
for installation instructions.

### Installing

CTK uses the Go mod system and is installed in any of the usual ways.

```
go get -u github.com/kckrinke/go-ctk
```

CTK includes a rudimentary clone of the venerable [dialog] application called
[go-dialog] and can be installed with the following:

```
go get github.com/kckrinke/go-ctk/cmd/go-dialog
```

CTK also includes a simple character-set viewer called [go-charmap] and can be
installed with the following:

```
go get github.com/kckrinke/go-ctk/cmd/go-charmap
```

## Example Usage

### Hello World in CTK

The following application will display an empty window, with "Hello World" as
the title and will quit when CTRL+C is pressed.

```
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
```

Compile the `hello-world.go` source file, and run it!

```
go build example/hello-world/hello-world.go && ./hello-world 
```

View the command-line help:

```
> ./hello-world --help
NAME:
   hello-world - the most basic CTK application

USAGE:
   hello-world [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cdk-log-file value   path to log file (default: /tmp/cdk.log) [$GO_CDK_LOG_FILE]
   --cdk-log-level value  highest level of verbosity (default: error) [$GO_CDK_LOG_LEVEL]
   --cdk-log-levels       list the levels of logging verbosity (default: false)
   --help, -h, --usage    display command-line usage information (default: false)
   --version              display the version (default: false)
```

### More Source Examples

For more complex examples, see the [go-dialog] and [go-charmap] sources, both of
which are intended as "complete project" examples to learn from and play with.

## Running the unit tests

```
go test -v
```

## Versioning

The current API is unstable and subject to change dramatically. The following is a brief summary of the planned iterations.

* v0.0.x - Proof of concept, experimental
* v0.1.x - Rewrite of CTK package, runtime systems
* v0.2.x - Use [CDK] v0.2.x, functional implementation
* v1.0.0 - First official release directly related to v1.0.0 of [CDK]

## Authors

* **Kevin C. Krinke** - *Original author* - [kckrinke]

## License

This project is licensed under the Apache License, Version 2.0 - see the
[LICENSE.md] file for details.

## Acknowledgments

* Thanks to [TCell] for providing a solid and robust platform to build upon
* Thanks to the [GTK Team] for developing and maintaining the [GTK2] API  that
  CTK is modeled after

[CDK]: https://github.com/kckrinke/go-cdk
[ctk-app.go]: example/ctk-app.go
[hello-world.go]: example/hello-world.go
[go-charmap]: cmd/go-charmap/main.go
[go-dialog]: cmd/go-dialog/main.go
[dialog]: https://invisible-island.net/dialog/
[kckrinke]: https://github.com/kckrinke
[LICENSE.md]: LICENSE.md
[TCell]: https://github.com/gdamore/tcell
[GTK Team]: https://www.gtk.org/development.php#Team
[GTK2]: https://developer.gnome.org/gtk2/
