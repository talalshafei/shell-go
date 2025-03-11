package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
	case "pwd":
		c.pwd()
	case "cd":
		c.cd()
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
	case "exit", "echo", "type", "pwd", "cd":
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

func (c *Command) pwd() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't retrieve the current working directory: %s", err.Error())
		return
	}
	fmt.Println(cwd)
}

func (c *Command) cd() {
	if len(c.Args) == 0 {
		return
	}

	if len(c.Args) > 1 {
		fmt.Println("bash: cd: too many arguments")
		return
	}

	newDir := c.Args[0]

	// handle '~'
	if newDir[0] == '~' {
		homeDir := os.Getenv("HOME")
		newDir = homeDir + newDir[1:]
	}

	err := os.Chdir(newDir)
	if err != nil {
		fmt.Printf("bash: cd: %s: No such file or directory\n", newDir)
	}
}

func (c *Command) run() {
	program := exec.Command(c.Name, c.Args...)
	program.Stdin = os.Stdin
	program.Stdout = os.Stdout
	program.Stderr = os.Stderr

	program.Run()
}

// Using built-in in GO
func (c *Command) searchPath() string {
	filePath, err := exec.LookPath(c.Name)
	if err != nil {
		return ""
	}
	return filePath
}
