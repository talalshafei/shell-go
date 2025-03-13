package main

import "github.com/codecrafters-io/shell-starter-go/app/shellparser"

func main() {
	parser := shellparser.NewParser()
	shell := NewShell(parser)

	shell.Start()
}
