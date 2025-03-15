package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/shellparser"
	"github.com/codecrafters-io/shell-starter-go/app/user_scanner"
)

func main() {
	uScanner := user_scanner.NewUserScanner(os.Stdin)
	parser := shellparser.NewParser()

	shell := NewShell(uScanner, parser)
	shell.Start()
}
