package editor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"unicode"
)

func _CTRL_KEY(char int) int {
	return (char) & 0x1f
}

const (
	BACKSPACE  = 127
	ARROW_LEFT = iota + 1000
	ARROW_RIGHT
	ARROW_UP
	ARROW_DOWN
	DEL_KEY
	HOME_KEY
	END_KEY
	PAGE_UP
	PAGE_DOWN
)

const PS1 = "$ "
const cursorStart = len(PS1) + 1

type Editor struct {
	*autoComplete
	*config
	cursor int
	Input  []byte
	rbuf   *bufio.Reader
}

func NewEditor() *Editor {
	c := &config{}
	c.enableRawMode()
	reader := bufio.NewReader(os.Stdin)
	ac := newAutoComplete()

	return &Editor{
		autoComplete: ac,
		config:       c,
		cursor:       cursorStart,
		Input:        nil,
		rbuf:         reader,
	}
}

func (e *Editor) cleanEditor() {
	e.Input = nil
	e.cursor = cursorStart
}

func (e *Editor) TakeInput() []byte {
	defer e.cleanEditor()
	fmt.Print(PS1)
	for e.processKeyPress() {
		e.refreshLine()
	}
	fmt.Println()
	return e.Input
}

func (e *Editor) Destroy() {
	e.config.disableRawMode()
}

func (e *Editor) processKeyPress() bool {
	c := e.readKey()

	switch c {
	case '\r', '\n':
		e.Input = append(e.Input, '\n') // for parser
		return false

	case '\t':
		e.handleAutoComplete()

	case _CTRL_KEY('l'):
		fmt.Print("\x1b[2J\x1b[H")

	case BACKSPACE, DEL_KEY, _CTRL_KEY('h'):
		if c == DEL_KEY {
			if e.cursor < cursorStart+len(e.Input) {
				e.moveCursor(ARROW_RIGHT)
				e.removeChar()
			}
		} else {
			e.removeChar()
		}

	case _CTRL_KEY('c'), _CTRL_KEY('d'), _CTRL_KEY('z'):
		// TODO should send signal to the exit function
		// for now do it yourself
		fmt.Println("^Interrupt")
		e.Destroy()
		os.Exit(0)

	case ARROW_LEFT, ARROW_RIGHT, ARROW_UP, ARROW_DOWN:
		e.moveCursor(c)

	default:
		if char := rune(c); unicode.IsPrint(char) {
			e.insertChar(byte(char))
		}
	}

	return true
}

func (e *Editor) refreshLine() {
	buf := []byte{}

	// "\x1b[?25l" hide cursor
	buf = append(buf, []byte("\x1b[?25l")...)
	// position to start of line
	buf = append(buf, []byte("\x1b[G")...)
	// "\x1b[K" everything right to the cursor
	buf = append(buf, []byte("\x1b[K")...)

	// "$ input" add data
	data := fmt.Sprintf("%s%s", PS1, e.Input)
	buf = append(buf, []byte(data)...) // might change

	// position cursor to the end of text
	position := fmt.Sprintf("\x1b[%dG", e.cursor)
	buf = append(buf, []byte(position)...)

	// "\x1b[?25h" show cursor
	buf = append(buf, []byte("\x1b[?25h")...)

	nwrite, err := os.Stdout.Write(buf)

	e.panicOnErr("refreshScreen", err)

	if nwrite != len(buf) {
		e.panicOnErr("refreshScreen", fmt.Errorf("wrote %d, but buffer size %d", nwrite, len(buf)))
	}
}

func (e *Editor) readKey() int {
	var err error
	var char byte

	char, err = e.rbuf.ReadByte()
	e.panicOnErr("readKey", err)

	// if ESC
	if char == '\x1b' {
		seq := [3]byte{}
		seq[0], err = e.rbuf.ReadByte()

		if err != nil {
			return int(char)
		}
		seq[1], err = e.rbuf.ReadByte()
		if err != nil {
			return int(char)
		}

		if seq[0] == '[' {
			if seq[1] >= '0' && seq[1] <= '9' {
				seq[2], err = e.rbuf.ReadByte()
				if err != nil {
					return int(char)
				}
				if seq[2] == '~' {
					switch seq[1] {
					case '1':
						return HOME_KEY
					case '3':
						return DEL_KEY
					case '4':
						return END_KEY
					case '5':
						return PAGE_UP
					case '6':
						return PAGE_DOWN
					case '7':
						return HOME_KEY
					case '8':
						return END_KEY
					}
				}
			} else {
				switch seq[1] {
				case 'A':
					return ARROW_UP
				case 'B':
					return ARROW_DOWN
				case 'C':
					return ARROW_RIGHT
				case 'D':
					return ARROW_LEFT
				case 'H':
					return HOME_KEY
				case 'F':
					return END_KEY
				}
			}
		} else if seq[0] == 'O' {
			switch seq[1] {
			case 'H':
				return HOME_KEY
			case 'F':
				return END_KEY
			}
		}

		return int(char)
	}
	return int(char)
}

func (e *Editor) insertChar(char byte) {
	at := e.cursor - cursorStart

	e.Input = append(e.Input, 0)
	copy(e.Input[at+1:], e.Input[at:])
	e.Input[at] = char

	e.cursor++
}

func (e *Editor) moveCursor(arrow int) {
	switch arrow {
	case ARROW_LEFT:
		if e.cursor > cursorStart {
			e.cursor--
		}
	case ARROW_RIGHT:
		if e.cursor < len(e.Input)+cursorStart {
			e.cursor++
		}
	case ARROW_UP, ARROW_DOWN:
		fmt.Print("\a")
		// later can implement history
	}
}

func (e *Editor) removeChar() {
	at := e.cursor - cursorStart - 1
	if at < 0 || at >= len(e.Input) {
		return
	}
	e.Input = slices.Delete(e.Input, at, at+1)
	e.cursor--
}

func (e *Editor) handleAutoComplete() {
	if len(e.Input) == 0 {
		fmt.Printf("\a")
		return
	}

	restOfWord, flag := e.autoComplete.completeWord(string(e.Input))
	if flag != FOUND_ONE {
		fmt.Printf("\a")
		return
	}
	e.Input = append(e.Input, []byte(restOfWord)...)
	e.Input = append(e.Input, ' ')
	e.cursor += len(restOfWord) + 1
}

func (e *Editor) panicOnErr(msg string, err error) {
	if err != nil {
		e.Destroy()
		log.Panicf("Editor: %s: %s\n", msg, err.Error())
	}
}
