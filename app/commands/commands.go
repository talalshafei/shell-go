package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Command struct {
	Name   string
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func StartCommand(input []string) {
	if len(input) == 0 {
		return
	}
	name := input[0]

	args, stdin, stdout, stderr, err := Redirect(input[1:])
	if err != nil {
		return
	}

	if stdin == nil {
		stdin = os.Stdin
	} else {
		defer stdin.Close()
	}

	if stdout == nil {
		stdout = os.Stdout
	} else {
		defer stdout.Close()
	}

	if stderr == nil {
		stderr = os.Stderr
	} else {
		defer stderr.Close()
	}

	cmd := &Command{
		Name:   name,
		Args:   args,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	cmd.Execute()
}

func NewCommand(name string, args []string) *Command {

	return &Command{
		Name:   name,
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
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
			fmt.Fprintf(c.Stderr, "%s: command not found\n", strings.Join(append([]string{c.Name}, c.Args...), " "))
		}
	}

}

func (c *Command) exit() {
	if len(c.Args) == 0 {
		fmt.Fprintln(c.Stderr, "Invalid exit code")
		return
	}
	code, err := strconv.Atoi(c.Args[0])
	if err != nil {
		fmt.Fprintln(c.Stderr, "Invalid exit code")
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
		fmt.Fprintf(c.Stdout, "%s is a shell builtin\n", typeCmd.Name)
	default:
		// executables found in PATH
		location := typeCmd.searchPath()
		if location != "" {
			fmt.Fprintf(c.Stdout, "%s is %s\n", typeCmd.Name, location)
		} else {
			fmt.Fprintf(c.Stderr, "%s: not found\n", typeCmd.Name)
		}
	}
}

func (c *Command) pwd() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(c.Stderr, "Couldn't retrieve the current working directory: %s", err.Error())
		return
	}
	fmt.Println(cwd)
}

func (c *Command) cd() {
	if len(c.Args) == 0 {
		return
	}

	if len(c.Args) > 1 {
		fmt.Fprintln(c.Stderr, "bash: cd: too many arguments")
		return
	}

	newDir := c.Args[0]
	err := os.Chdir(newDir)

	if err != nil {
		fmt.Fprintf(c.Stderr, "bash: cd: %s: No such file or directory\n", newDir)
	}
}

func (c *Command) run() {
	program := exec.Command(c.Name, c.Args...)
	program.Stdin = c.Stdin
	program.Stdout = c.Stdout
	program.Stderr = c.Stderr

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
