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

		if cmd == "exit" {
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
		} else if cmd == "" {
			continue

		} else {
			fmt.Printf("%s: command not found\n", input)
		}
	}
}
