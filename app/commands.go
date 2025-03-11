package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func TakeCommand() *Command {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading the input: ", err.Error())
		return NewCommand("", nil)
	}
	input = strings.TrimRight(input, "\r\n")

	elems := strings.Split(input, " ")

	return NewCommand(elems[0], elems[1:])
}
func NewCommand(name string, args []string) *Command {
	return &Command{name, args}
}

func (c *Command) Execute() {

	// builtin and empty string
	switch c.Name {
	case "":
		return
	case "exit":
		c.exit()
	case "echo":
		c.echo()
	case "type":
		c.typeCommand()
	default:
		// executables found in PATH
		location := c.searchPath()
		if location != "" {
			c.run()
		} else {
			fmt.Printf("%s: command not found\n", strings.Join(append([]string{c.Name}, c.Args...), " "))
		}
	}

}

func (c *Command) exit() {
	if len(c.Args) == 0 {
		fmt.Println("Invalid exit code")
		return
	}
	code, err := strconv.Atoi(c.Args[0])
	if err != nil {
		fmt.Println("Invalid exit code")
		return
	}
	os.Exit(code)
}

func (c *Command) echo() {
	fmt.Println(strings.Join(c.Args, " "))
}

func (c *Command) typeCommand() {
	var name string
	if len(c.Args) != 0 {
		name = c.Args[0]
	}

	typeCmd := NewCommand(name, nil)
	switch typeCmd.Name {
	case "exit", "echo", "type":
		fmt.Printf("%s is a shell builtin\n", typeCmd.Name)
	default:
		// executables found in PATH
		location := typeCmd.searchPath()
		if location != "" {
			fmt.Printf("%s is %s\n", typeCmd.Name, location)
		} else {
			fmt.Printf("%s: not found\n", typeCmd.Name)
		}
	}
}

func (c *Command) run() {

}

// Correct way manually
func (c *Command) searchPath() string {
	pathENV := os.Getenv("PATH")

	for dir := range strings.SplitSeq(pathENV, ":") {
		fullPath := filepath.Join(dir, c.Name)
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			if info.Mode()&0111 != 0 {
				return fullPath
			}
		}
	}
	return ""
}
