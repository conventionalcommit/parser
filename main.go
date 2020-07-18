package main

import (
	"os"

	"github.com/xfxdev/xlog"
)

// applicationVersion is set by the linker
var applicationVersion = "0.0.0"

// main is the entrypoint when called from the command line
func main() {
	setup()
	if err := rootCmd.Execute(); err != nil {
		xlog.Error(err)
		os.Exit(1)
	}
}
