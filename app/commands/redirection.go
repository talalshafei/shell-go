package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrUnexpectedTokenRedirect = errors.New("bash: syntax error near unexpected token `newline'")

func Redirect(input []string) (args []string, stdin *os.File, stdout *os.File, stderr *os.File, err error) {

	for i := range input {
		token := input[i]
		remaining := input[i+1:]

		if isRedirectOperator(token) {
			switch token {
			case "1>":
				stdout, err = prepareOutput(remaining, false)
			case "2>":
				stderr, err = prepareOutput(remaining, false)
			case "1>>":
				stdout, err = prepareOutput(remaining, true)
			case "2>>":
				stderr, err = prepareOutput(remaining, true)
			case "0<":
				stdin, err = prepareInput(remaining)
			}
			break
		} else {
			args = append(args, token)
		}
	}

	if err != nil {
		return nil, nil, nil, nil, err
	}

	return args, stdin, stdout, stderr, nil
}

func isRedirectOperator(token string) bool {
	if token == "1>" || token == "2>" || token == "1>>" || token == "2>>" || token == "0<" {
		return true
	}
	return false
}

func prepareOutput(path []string, append bool) (*os.File, error) {
	if len(path) == 0 || path[0] == "\n" {
		return nil, ErrUnexpectedTokenRedirect
	}
	fmt.Printf("%v %t\n", path, append)
	filepathStr := strings.Join(path, "")
	dirStr := filepath.Dir(filepathStr)

	if err := os.MkdirAll(dirStr, 0777); err != nil {
		return nil, err
	}

	flags := os.O_WRONLY | os.O_CREATE
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(filepathStr, flags, 0777)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func prepareInput(path []string) (*os.File, error) {
	if len(path) == 0 || path[0] == "\n" {
		return nil, ErrUnexpectedTokenRedirect
	}

	filepathStr := strings.Join(path, "")
	file, err := os.Open(filepathStr)
	if err != nil {
		return nil, err
	}

	return file, nil
}
