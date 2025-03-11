package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
		fmt.Printf("%s is a shell builtin\n", typeCmd)
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

// Note recursive search is not really needed
// Implementing it the hard way just for fun
func (c *Command) searchPath() string {
	pathENV := os.Getenv("PATH")

	for dir := range strings.SplitSeq(pathENV, ":") {
		if location, err := findFileRecursively(dir, c.Name); err == nil && location != "" {
			return location
		}
	}
	return ""
}

func findFileRecursively(dir, exeName string) (string, error) {
	d, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer d.Close()

	entries, err := d.Readdir(-1)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		fullPath := joinPaths(dir, entry.Name())

		if entry.IsDir() {
			if location, err := findFileRecursively(fullPath, exeName); err == nil && location != "" {
				return location, nil
			}
		} else {
			if entry.Name() == exeName && (entry.Mode()&0111) != 0 {
				return fullPath, nil
			}
		}
	}

	return "", errors.New("not found")
}

func joinPaths(dir, entryName string) string {
	if strings.HasSuffix(dir, "/") {
		return dir + entryName
	}
	return dir + "/" + entryName
}
