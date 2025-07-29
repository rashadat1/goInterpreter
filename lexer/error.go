package lexer

import (
	"fmt"
	"strconv"
)

type TokenError interface {
	FormatMessage() string
}

type UnrecognizedCharError struct {
	Message string
	Char    string
	Line    int
}
type UnterminatedStringError struct {
	Message string
	Line    int
}

func (uc UnrecognizedCharError) FormatMessage() string {
	return fmt.Sprintf("[line %s] Error: Unexpected character: %s\n", strconv.Itoa(uc.Line), uc.Char)
}

func (us UnterminatedStringError) FormatMessage() string {
	return fmt.Sprintf("[line %s] Error: Unterminated string.", strconv.Itoa(us.Line))
}

type TokenErrors []TokenError

func (te TokenErrors) ToString() string {
	out := ""
	for _, tokErr := range te {
		out += tokErr.FormatMessage()
	}
	return out
}
