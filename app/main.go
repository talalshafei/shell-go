package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		fmt.Print("$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the input: ", err.Error())
			return
		}

		input = strings.TrimRight(input, "\r\n")

		args := strings.Split(input, " ")

		cmd := args[0]

		switch cmd {
		case "":
			continue
		case "exit":
			if len(args) != 2 {
				fmt.Println("Invalid exit code")
				continue
			}
			code, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid exit code")
				continue
			}
			os.Exit(code)
		case "echo":
			fmt.Println(strings.Join(args[1:], " "))
		case "type":
			var typeCmd string
			if len(args) > 1 {
				typeCmd = args[1]
			}
			switch typeCmd {
			case "exit", "echo", "type":
				fmt.Printf("%s is a shell builtin\n", typeCmd)
			default:
				fmt.Printf("%s: command not found\n", typeCmd)
			}

		default:
			fmt.Printf("%s: command not found\n", input)
		}

	}
}
