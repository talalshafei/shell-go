package main

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"github.com/codecrafters-io/shell-starter-go/app/editor"
	"github.com/codecrafters-io/shell-starter-go/app/shellparser"
)

type Shell struct {
	editor *editor.Editor
	parser *shellparser.Parser
}

func NewShell(e *editor.Editor, p *shellparser.Parser) *Shell {
	return &Shell{
		editor: e,
		parser: p,
	}
}

func (sh *Shell) Start() int {
	isExit, exitCode := false, 0

	editor := sh.editor
	parser := sh.parser

	for !isExit {
		var err error

		// take input
		rawInput := editor.TakeInput()

		// parse input into
		inputStringArr, err := parser.Parse(rawInput)

		if err != nil {
			fmt.Println(err)
			continue
		}

		// prepare and start commands
		isExit, exitCode = commands.StartCommand(inputStringArr)
	}

	return exitCode
}
