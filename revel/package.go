package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apinprastya/revel"
)

var cmdPackage = &Command{
	UsageLine: "package [import path] [run mode]",
	Short:     "package a Revel application (e.g. for deployment)",
	Long: `
Package the Revel web application named by the given import path.
This allows it to be deployed and run on a machine that lacks a Go installation.

The run mode is used to select which set of app.conf configuration should
apply and may be used to determine logic in the application itself.

Run mode defaults to "dev".

For example:

    revel package github.com/apinprastya/samples/chat
`,
}

func init() {
	cmdPackage.Run = packageApp
}

func packageApp(args []string) {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, cmdPackage.Long)
		return
	}

	// Determine the run mode.
	mode := "dev"
	if len(args) >= 2 {
		mode = args[1]
	}

	appImportPath := args[0]
	revel.Init(mode, appImportPath, "")

	// Remove the archive if it already exists.
	destFile := filepath.Base(revel.BasePath) + ".tar.gz"
	os.Remove(destFile)

	// Collect stuff in a temp directory.
	tmpDir, err := ioutil.TempDir("", filepath.Base(revel.BasePath))
	panicOnError(err, "Failed to get temp dir")

	buildApp([]string{args[0], tmpDir, mode})

	// Create the zip file.
	archiveName := mustTarGzDir(destFile, tmpDir)

	fmt.Println("Your archive is ready:", archiveName)
}
