package main

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/app/shellparser"
)

type Shell struct {
	parser *shellparser.Parser
}

func NewShell(parser *shellparser.Parser) *Shell {
	return &Shell{parser}
}

func (sh *Shell) Start() {
	parser := sh.parser
	for {
		fmt.Print("$ ")
		input, err := parser.TakeAndParseInput()
		if err != nil {
			fmt.Println(err)
			continue
		}
		StartCommand(input)
	}
}
