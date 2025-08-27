package parser

import (
	"fmt"
	"strconv"

	"github/goInterpreter/lexer"
)

type ParserError struct {
	Line    int
	Message string
}

func (pe ParserError) Report(token lexer.Token) string {
	var errorMsg string
	if token.Type == lexer.TokenEOF {
		errorMsg = fmt.Sprintf("[line %s] Error at end: %s", strconv.Itoa(pe.Line), pe.Message)
	} else {
		errorMsg = fmt.Sprintf("[line %s] Error at %s: %s", strconv.Itoa(pe.Line), token.Lexeme, pe.Message)
	}
	return errorMsg
}
