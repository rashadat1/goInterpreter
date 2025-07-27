package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenLeftParen
	TokenRightParen
	TokenNewLine
)

type TokenizedText []Token

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
	//lexeme bytes.Buffer
}

// Emit a Token by reading the next rune from the bytes.Reader object stored in Lexer
func (l *Lexer) NextToken() (Token, error) {
	r := l.reader
	c, _, err := r.ReadRune()
	if err != nil {
		return Token{}, err
	}
	switch c {
	case '(':
		return Token{
			Type:    TokenLeftParen,
			Lexeme:  "(",
			Literal: "null",
			Line:    l.line,
		}, nil
	case ')':
		return Token{
			Type:    TokenRightParen,
			Lexeme:  ")",
			Literal: "null",
			Line:    l.line,
		}, nil
	case '\n':
		l.line += 1
		return Token{
			Type: TokenNewLine,
		}, nil
	default:
		err := errors.New("Token not implemented: " + string(c))
		return Token{}, err
	}
}

func (l *Lexer) ScanTokens() (TokenizedText, error) {
	toks := make(TokenizedText, 0)
	for {
		tok, err := l.NextToken()
		if err != nil {
			if err == io.EOF {
				return toks, nil
			}
			return nil, err
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
		run(lex)
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
		run(lex)
		fmt.Printf("Input code: %s\n", line)
	}
}

func run(r Lexer) {
	toks, err := r.ScanTokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred processing tokens: %s", err.Error())
	}
	res := toks.toString()
	fmt.Print(res)
}
