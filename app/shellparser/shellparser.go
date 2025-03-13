package shellparser

import (
	"bufio"
	"errors"
	"os"
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

func (p *Parser) TakeAndParseInput() ([]string, error) {
	input, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return p.Parse(input)
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

	p.finalizeCurrentToken()
	return p.tokens, nil
}

func (p *Parser) handleNormalState(char byte, input []byte, idx *int) {
	switch {
	case unicode.IsSpace(rune(char)):
		p.handleWhitespace()
	case char == '\'':
		p.state = stateSingleQuote
	case char == '"':
		p.state = stateDoubleQuote
	case char == '\\':
		p.handleBackslashEscape(input, idx)
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

func (p *Parser) finalizeCurrentToken() {
	if p.currentToken.Len() > 0 {
		p.tokens = append(p.tokens, p.currentToken.String())
		p.currentToken.Reset()
	}
}
