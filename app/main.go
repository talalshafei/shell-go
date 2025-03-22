package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/editor"
	"github.com/codecrafters-io/shell-starter-go/app/shellparser"
)

// register an exit function

func main() {
	editor := editor.NewEditor()
	parser := shellparser.NewParser()
	shell := NewShell(editor, parser)
	exitCode := shell.Start()

	editor.Destroy()
	os.Exit(exitCode)
}
