package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading the input: ", err.Error())
		return
	}

	input = strings.TrimRight(input, "\r\n")
	fmt.Printf("%s: command not found\n", input)
}
