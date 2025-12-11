package main

import (
	"fmt"
	"os"

	"github.com/idongju/t9s/internal/ui"
)

const (
	appVersion = "1.0.0"
	appName    = "T9s"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	// Use new architecture (v0.2.0)
	// To use legacy app: ui.NewApp()
	app := ui.NewAppNew()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
