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
	filepathStr := strings.Join(path, "")

	// Get the directory part
	dirStr := filepath.Dir(filepathStr)

	// Create the directory with specific error handling
	if dirStr != "." {
		err := os.MkdirAll(dirStr, 0777)
		if err != nil {
			// Log the specific error for debugging
			fmt.Fprintf(os.Stderr, "Failed to create directory '%s': %v\n", dirStr, err)
			return nil, err
		}
	}

	// Verify the directory exists before proceeding
	dirInfo, err := os.Stat(dirStr)
	if err != nil || !dirInfo.IsDir() {
		return nil, fmt.Errorf("directory '%s' could not be created or accessed", dirStr)
	}

	flags := os.O_WRONLY | os.O_CREATE
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(filepathStr, flags, 0666)

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
