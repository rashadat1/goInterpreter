package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

type TokenType int

type TokenError interface {
	FormatMessage() string
}
type UnrecognizedCharError struct {
	Message string
	Char    string
	Line    int
}

func (uc UnrecognizedCharError) FormatMessage() string {
	return fmt.Sprintf("[line %s] Error: Unexpected character: %s\n", strconv.Itoa(uc.Line), uc.Char)
}

const (
	TokenEOF TokenType = iota
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenNewLine
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemiColon
	TokenSlash
	TokenStar
)

type TokenizedText []Token
type TokenErrors []TokenError

func (t TokenType) toString() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenLeftParen:
		return "LEFT_PAREN"
	case TokenRightParen:
		return "RIGHT_PAREN"
	case TokenNewLine:
		return "NEW_LINE"
	case TokenLeftBrace:
		return "LEFT_BRACE"
	case TokenRightBrace:
		return "RIGHT_BRACE"
	case TokenComma:
		return "COMMA"
	case TokenDot:
		return "DOT"
	case TokenMinus:
		return "MINUS"
	case TokenPlus:
		return "PLUS"
	case TokenSemiColon:
		return "SEMICOLON"
	case TokenSlash:
		return "SLASH"
	case TokenStar:
		return "STAR"
	default:
		return ""
	}
}

func (tt TokenizedText) toString() string {
	out := ""
	for _, tok := range tt {
		out += tok.TokToString()
	}
	eofTok := Token{
		Type:    TokenEOF,
		Lexeme:  "",
		Literal: "null",
	}
	out += eofTok.TokToString()
	return out
}

func (te TokenErrors) toString() string {
	out := ""
	for _, tokErr := range te {
		out += tokErr.FormatMessage()
	}
	return out
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func LiteralToString(lit interface{}) string {
	switch v := lit.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func (t *Token) TokToString() string {
	return t.Type.toString() + " " + t.Lexeme + " " + LiteralToString(t.Literal) + string('\n')
}

type Lexer struct {
	reader *bytes.Reader
	line   int
	lexErr TokenErrors
	//lexeme bytes.Buffer
}

// Emit a Token by reading the next rune from the bytes.Reader object stored in Lexer
func (l *Lexer) NextToken() (Token, TokenError, error) {
	r := l.reader
	c, _, err := r.ReadRune()
	if err != nil {
		return Token{}, nil, err
	}
	switch c {
	// single-char tokens
	case '(':
		return Token{
			Type:    TokenLeftParen,
			Lexeme:  "(",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case ')':
		return Token{
			Type:    TokenRightParen,
			Lexeme:  ")",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case '{':
		return Token{
			Type:    TokenLeftBrace,
			Lexeme:  "{",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case '}':
		return Token{
			Type:    TokenRightBrace,
			Lexeme:  "}",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case ',':
		return Token{
			Type:    TokenComma,
			Lexeme:  ",",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case '.':
		return Token{
			Type:    TokenDot,
			Lexeme:  ".",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case '-':
		return Token{
			Type:    TokenMinus,
			Lexeme:  "-",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case '+':
		return Token{
			Type:    TokenPlus,
			Lexeme:  "+",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case ';':
		return Token{
			Type:    TokenSemiColon,
			Lexeme:  ";",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case '/':
		return Token{
			Type:    TokenSlash,
			Lexeme:  "/",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case '*':
		return Token{
			Type:    TokenStar,
			Lexeme:  "*",
			Literal: "null",
			Line:    l.line,
		}, nil, nil
	case '\n':
		l.line += 1
		return Token{
			Type: TokenNewLine,
		}, nil, nil
	default:
		uc := UnrecognizedCharError{
			Char: string(c),
			Line: l.line,
		}
		uc.Message = uc.FormatMessage()
		return Token{}, uc, nil
	}
}

func (l *Lexer) ScanTokens() (TokenizedText, TokenErrors, error) {
	toks := make(TokenizedText, 0)
	tokErrs := make(TokenErrors, 0)
	for {
		tok, tokErr, err := l.NextToken()
		if err != nil {
			if err == io.EOF {
				return toks, tokErrs, nil
			}
			return nil, nil, err
		}
		if tokErr != nil {
			tokErrs = append(tokErrs, tokErr)
			continue
		}
		if tok.Type.toString() != "NEW_LINE" {
			toks = append(toks, tok)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: ./your_program.sh tokenize <filename>\n")
		fmt.Fprintf(os.Stderr, "Usage: ./your_program.sh tokenize")
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		replEcho()
		os.Exit(0)
	}
	command := os.Args[1]
	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
	fileName := os.Args[2]
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %s", fileName, err)
		os.Exit(1)
	}
	if len(fileContents) > 0 {
		r := bytes.NewReader(fileContents)

		lex := Lexer{
			reader: r,
			line:   1,
		}
		run(&lex)
		if len(lex.lexErr) != 0 {
			os.Exit(65)
		}
		os.Exit(0)
	} else {
		fmt.Println("EOF  null")
	}
}

func replEcho() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprintf(os.Stdout, "> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error occurred reading line: %s", err.Error())
				continue
			}
			// if scanner.Scan() == false but not error occurred it means user hit control + D
			break
		}
		line := scanner.Text()
		if line == "" || line == string('\n') {
			continue
		}
		lineBytes := []byte(line)
		r := bytes.NewReader(lineBytes)
		lex := Lexer{
			reader: r,
			line:   1,
		}
		run(&lex)
		fmt.Printf("Input code: %s\n", line)
	}
}

func run(lex *Lexer) {
	res := ""
	resErr := ""
	toks, tokErrs, err := lex.ScanTokens()
	lex.lexErr = tokErrs
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred processing tokens: %s", err.Error())
	}
	res = toks.toString()
	resErr = tokErrs.toString()
	fmt.Print(res)
	fmt.Fprint(os.Stderr, resErr)
}
