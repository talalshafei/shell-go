package user_scanner

import (
	"bufio"
	"io"
)

type UserScanner struct {
	reader io.Reader
	Text   []byte
}

func NewUserScanner(file io.Reader) *UserScanner {
	return &UserScanner{
		reader: file,
	}
}

func (u *UserScanner) CaptureInput() error {
	var err error
	u.Text, err = bufio.NewReader(u.reader).ReadBytes('\n')
	return err
}
