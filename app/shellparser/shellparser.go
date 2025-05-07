package shellparser

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Parser struct {
	tokens       []string
	currentToken strings.Builder
	state        int
	err          error
}

const (
	stateNormal = iota
	stateSingleQuote
	stateDoubleQuote
)

// refactor later to be pause and open new line ">" instead of error to take more input from shell
var (
	ErrUnclosedQuotes    = errors.New("unclosed quotes")
	ErrBackslashAtEnd    = errors.New("backslash at end of input")
	ErrDanglingBackslash = errors.New("dangling backslash in double quotes")

	// real error
	ErrUnexpectedTokenRedirect = errors.New("bash: syntax error near unexpected token `newline'")
	ErrUnexpectedTokenPipe     = errors.New("bash: syntax error near unexpected token `|'")
)

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) cleanParser() {
	p.tokens = nil
	p.currentToken.Reset()
	p.state = stateNormal
	p.err = nil
}

func (p *Parser) Parse(input []byte) ([]string, error) {
	defer p.cleanParser()

	for i := 0; i < len(input) && p.err == nil; i++ {
		char := input[i]
		switch p.state {
		case stateNormal:
			p.handleNormalState(char, input, &i)
		case stateSingleQuote:
			p.handleSingleQuoteSate(char)
		case stateDoubleQuote:
			p.handleDoubleQuoteState(char, input, &i)
		}
	}

	if p.err != nil {
		return nil, p.err
	}

	if p.state != stateNormal {
		return nil, ErrUnclosedQuotes
	}

	p.flushCurrentToken()
	return p.tokens, nil
}

func (p *Parser) handleNormalState(char byte, input []byte, idx *int) {
	switch {
	case char == '>':
		p.handleOutputRedirect(input, idx)
	case char == '<':
		p.handleInputRedirect(input, idx)

	case unicode.IsSpace(rune(char)):
		p.handleWhitespace()
	case char == '\'':
		p.state = stateSingleQuote
	case char == '"':
		p.state = stateDoubleQuote

	case char == '\\':
		p.handleBackslashEscape(input, idx)

	case char == '|':
		if (len(p.tokens) == 0 && p.currentToken.Len() == 0) || (len(p.tokens) > 0 && p.tokens[len(p.tokens)-1] == "|") {
			p.err = ErrUnexpectedTokenPipe
			return
		}
		p.flushCurrentToken()
		p.tokens = append(p.tokens, "|")

	default:
		p.currentToken.WriteByte(char)
	}
}

func (p *Parser) handleSingleQuoteSate(char byte) {
	if char == '\'' {
		p.state = stateNormal
	} else {
		p.currentToken.WriteByte(char)
	}
}

func (p *Parser) handleDoubleQuoteState(char byte, input []byte, idx *int) {
	switch char {
	case '"':
		p.state = stateNormal
	case '\\':
		p.handleDoubleQuoteBackslash(input, idx)
	default:
		p.currentToken.WriteByte(char)
	}
}

func (p *Parser) handleWhitespace() {
	if p.currentToken.Len() > 0 {
		out := p.currentToken.String()
		if out[0] == '~' {
			out = os.Getenv("HOME") + out[1:]
		}
		p.tokens = append(p.tokens, out)
		p.currentToken.Reset()
	}
}

func (p *Parser) handleBackslashEscape(input []byte, idx *int) {
	if *idx+1 >= len(input) {
		p.err = ErrBackslashAtEnd
		return
	}

	nextChar := input[*idx+1]
	if nextChar == '\n' {
		*idx++
	} else {
		p.currentToken.WriteByte(nextChar)
		*idx++
	}
}

func (p *Parser) handleDoubleQuoteBackslash(input []byte, idx *int) {
	if *idx+1 >= len(input) {
		p.err = ErrDanglingBackslash
		return
	}
	nextChar := input[*idx+1]
	switch nextChar {
	case '"', '$', '`', '\\':
		p.currentToken.WriteByte(nextChar)
		*idx++
	case '\n':
		*idx++
	default:
		p.currentToken.WriteByte('\\')
		p.currentToken.WriteByte(nextChar)
		*idx++
	}
}

func (p *Parser) flushCurrentToken() {
	if p.currentToken.Len() > 0 {
		p.tokens = append(p.tokens, p.currentToken.String())
		p.currentToken.Reset()
	}
}

func isStringFileDescriptor(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// > or >>
func (p *Parser) handleOutputRedirect(input []byte, idx *int) {

	if p.raiseUnexpectedTokenRedirection(len(input), *idx) {
		return
	}

	redirectString := "1>"

	cur := p.currentToken.String()
	if isStringFileDescriptor(cur) {
		if cur == "2" {
			redirectString = "2>"
		}
		p.currentToken.Reset()
	} else {
		p.flushCurrentToken()
	}

	// append
	if input[*idx+1] == '>' {
		redirectString += ">"
		*idx++
		if p.raiseUnexpectedTokenRedirection(len(input), *idx) {
			return
		}
	}

	p.tokens = append(p.tokens, redirectString)
}

// <
func (p *Parser) handleInputRedirect(input []byte, idx *int) {
	if p.raiseUnexpectedTokenRedirection(len(input), *idx) {
		return
	}
	if isStringFileDescriptor(p.currentToken.String()) {
		// remove file descriptor
		p.currentToken.Reset()
	} else {
		p.flushCurrentToken()
	}
	p.tokens = append(p.tokens, "0<")
}

func (p *Parser) raiseUnexpectedTokenRedirection(inputSize, idx int) bool {
	if idx+1 >= inputSize {
		p.err = ErrUnexpectedTokenRedirect
		return true
	}
	return false
}
