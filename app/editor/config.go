package editor

import (
	"log"
	"os"

	"golang.org/x/sys/unix"
)

type config struct {
	oldState *unix.Termios
}

func (c *config) enableRawMode() {
	var err error

	fd := int(os.Stdin.Fd())
	// raw mode the hard way
	c.oldState, err = unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		log.Panicf("enableRawMode: %s\r\n", err.Error())
	}

	oldState := *c.oldState
	raw := &oldState

	// ICRNL flag of replacing '\r' with '\n'
	// IXON flag for CTRL+S and CTRL+Q
	// OPOST flag to for post processing
	// ECHO to print typed chars
	// ICANON flag for Canonical mode
	// IEXTEN for CTRL+V that lets you print control chars
	// ISIG for CTRL+C and CTRL+Z signals

	raw.Iflag &^= unix.ICRNL | unix.IXON
	// raw.Oflag &^= unix.OPOST // no need for disabling this
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.IEXTEN | unix.ISIG

	// the following flags may be already
	// turned off on modern shells but
	// we turn them off for completeness
	raw.Iflag &^= unix.BRKINT | unix.INPCK | unix.ISTRIP

	// this not a flag this sets the char size to 8-bits
	raw.Cflag |= unix.CS8

	// Timeouts
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0 // blocking until 1 byte

	if err = unix.IoctlSetTermios(fd, unix.TCSETS, raw); err != nil {
		log.Panicf("enableRawMode: %s\r\n", err.Error())
	}

}

func (c *config) disableRawMode() {
	if err := unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETS, c.oldState); err != nil {
		log.Panicf("enableRawMode: %s\r\n", err.Error())
	}
}
