package main

import (
	"fmt"
	"os"

	"github.com/idongju/t9s/internal/ui"
)

const (
	appVersion = "0.1.0"
	appName    = "T9s"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	app := ui.NewApp()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
