package main

import (
	"fmt"
)

type Shell struct{}

func NewShell() *Shell {
	return &Shell{}
}

func (sh *Shell) Start() {
	for {
		fmt.Print("$ ")
		TakeCommand().Execute()
	}
}
