package main

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"github.com/codecrafters-io/shell-starter-go/app/shellparser"
	"github.com/codecrafters-io/shell-starter-go/app/user_scanner"
)

type Shell struct {
	uScanner *user_scanner.UserScanner
	parser   *shellparser.Parser
}

func NewShell(u *user_scanner.UserScanner, p *shellparser.Parser) *Shell {
	return &Shell{
		uScanner: u,
		parser:   p,
	}
}

func (sh *Shell) Start() {
	uScanner := sh.uScanner
	parser := sh.parser

	for {
		fmt.Print("$ ")

		var err error
		// take input
		err = uScanner.CaptureInput()
		if err != nil {
			fmt.Println(err)
			continue
		}
		rawInput := uScanner.Text

		// parse input into
		inputStringArr, err := parser.Parse(rawInput)

		if err != nil {
			fmt.Println(err)
			continue
		}

		// prepare and start commands
		commands.StartCommand(inputStringArr)
	}

}
