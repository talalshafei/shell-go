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

func StartCommand(input []string) (bool, int) {
	if len(input) == 0 {
		return false, 0
	}
	name := input[0]

	args, stdin, stdout, stderr, err := Redirect(input[1:])

	if err != nil {
		fmt.Println(err)
		return false, 0
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

	return cmd.Execute()
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

// exitCode can be used later to report the exit code of every command
func (c *Command) Execute() (isExit bool, exitCode int) {

	// builtin and empty string
	switch c.Name {
	case "":
		return
	case "exit":
		isExit, exitCode = c.exit()
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

	return isExit, exitCode
}

func (c *Command) exit() (bool, int) {
	if len(c.Args) == 0 {
		fmt.Fprint(c.Stderr, "Invalid exit code\n")
		return false, 0
	}
	code, err := strconv.Atoi(c.Args[0])
	if err != nil {
		fmt.Fprint(c.Stderr, "Invalid exit code\n")
		return false, 0
	}
	return true, code
}

func (c *Command) echo() {
	fmt.Fprintf(c.Stdout, "%s\n", strings.Join(c.Args, " "))
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
		fmt.Fprintf(c.Stderr, "Couldn't retrieve the current working directory: %s\n", err.Error())
		return
	}
	fmt.Fprintf(c.Stdout, "%s\n", cwd)
}

func (c *Command) cd() {
	if len(c.Args) == 0 {
		return
	}

	if len(c.Args) > 1 {
		fmt.Fprint(c.Stderr, "bash: cd: too many arguments\n")
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
