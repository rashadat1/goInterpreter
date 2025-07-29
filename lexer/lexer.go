package lexer

import (
	"bytes"
	"io"
)

type Lexer struct {
	Reader *bytes.Reader
	Line   int
	//lexeme bytes.Buffer
}

// Emit a Token by reading the next rune from the bytes.Reader object stored in Lexer
func (l *Lexer) NextToken() (Token, TokenError, error) {
	r := l.Reader
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
			Line:    l.Line,
		}, nil, nil
	case ')':
		return Token{
			Type:    TokenRightParen,
			Lexeme:  ")",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case '{':
		return Token{
			Type:    TokenLeftBrace,
			Lexeme:  "{",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case '}':
		return Token{
			Type:    TokenRightBrace,
			Lexeme:  "}",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case ',':
		return Token{
			Type:    TokenComma,
			Lexeme:  ",",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case '.':
		return Token{
			Type:    TokenDot,
			Lexeme:  ".",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case '-':
		return Token{
			Type:    TokenMinus,
			Lexeme:  "-",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case '+':
		return Token{
			Type:    TokenPlus,
			Lexeme:  "+",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case '=':
		cn, _, err := r.ReadRune()
		if err == io.EOF {
			return Token{
				Type:    TokenEqual,
				Lexeme:  "=",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
		if err != nil {
			return Token{}, nil, err
		}
		switch cn {
		case '=':
			return Token{
				Type:    TokenEqualEqual,
				Lexeme:  "==",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		default:
			err := r.UnreadRune()
			if err != nil {
				return Token{}, nil, err
			}
			return Token{
				Type:    TokenEqual,
				Lexeme:  "=",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
	case '!':
		cn, _, err := r.ReadRune()
		if err == io.EOF {
			return Token{
				Type:    TokenBang,
				Lexeme:  "!",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
		if err != nil {
			return Token{}, nil, err
		}
		switch cn {
		case '=':
			return Token{
				Type:    TokenBangEqual,
				Lexeme:  "!=",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		default:
			err := r.UnreadRune()
			if err != nil {
				return Token{}, nil, err
			}
			return Token{
				Type:    TokenBang,
				Lexeme:  "!",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
	case '<':
		cn, _, err := r.ReadRune()
		if err == io.EOF {
			return Token{
				Type:    TokenLess,
				Lexeme:  "<",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
		if err != nil {
			return Token{}, nil, err
		}
		switch cn {
		case '=':
			return Token{
				Type:    TokenLessEqual,
				Lexeme:  "<=",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		default:
			err := r.UnreadRune()
			if err != nil {
				return Token{}, nil, err
			}
			return Token{
				Type:    TokenLess,
				Lexeme:  "<",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
	case '>':
		cn, _, err := r.ReadRune()
		if err == io.EOF {
			return Token{
				Type:    TokenGreater,
				Lexeme:  ">",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
		if err != nil {
			return Token{}, nil, err
		}
		switch cn {
		case '=':
			return Token{
				Type:    TokenGreaterEqual,
				Lexeme:  ">=",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		default:
			err := r.UnreadRune()
			if err != nil {
				return Token{}, nil, err
			}
			return Token{
				Type:    TokenGreater,
				Lexeme:  ">",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
	case ';':
		return Token{
			Type:    TokenSemiColon,
			Lexeme:  ";",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case '/':
		cn, _, err := r.ReadRune()
		if err == io.EOF {
			return Token{
				Type:    TokenSlash,
				Lexeme:  "/",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}
		if err != nil {
			return Token{}, nil, err
		}
		switch cn {
		case '/':
			// single-line comment
			for {
				// consume until '\n' or EOF
				nxt, _, err := r.ReadRune()
				if nxt == '\n' {
					// if end of line
					l.Line++
					return l.NextToken()
				}
				if err != nil {
					// or err occurs including io.EOF -> don't yield any tokens because
					// we were in a comment state
					return Token{}, nil, err
				}
			}
		default:
			err := r.UnreadRune()
			if err != nil {
				return Token{}, nil, err
			}
			return Token{
				Type:    TokenSlash,
				Lexeme:  "/",
				Literal: "null",
				Line:    l.Line,
			}, nil, nil
		}

	case '*':
		return Token{
			Type:    TokenStar,
			Lexeme:  "*",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case '\n':
		l.Line += 1
		return Token{
			Type: TokenNewLine,
		}, nil, nil
	default:
		uc := UnrecognizedCharError{
			Char: string(c),
			Line: l.Line,
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
		if tok.Type.ToString() != "NEW_LINE" {
			toks = append(toks, tok)
		}
	}
}
