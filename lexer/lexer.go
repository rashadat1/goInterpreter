package lexer

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	Reader *bytes.Reader
	Line   int
	Lexeme *bytes.Buffer
}

// Emit a Token by reading the next rune from the bytes.Reader object stored in Lexer
func (l *Lexer) NextToken() (Token, TokenError, error) {
	r := l.Reader
	c, _, err := r.ReadRune()
	if err != nil {
		return Token{}, nil, err
	}
	switch {
	// single-char tokens
	case c == '(':
		return Token{
			Type:    TokenLeftParen,
			Lexeme:  "(",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == ')':
		return Token{
			Type:    TokenRightParen,
			Lexeme:  ")",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == '{':
		return Token{
			Type:    TokenLeftBrace,
			Lexeme:  "{",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == '}':
		return Token{
			Type:    TokenRightBrace,
			Lexeme:  "}",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == ',':
		return Token{
			Type:    TokenComma,
			Lexeme:  ",",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == '.':
		return Token{
			Type:    TokenDot,
			Lexeme:  ".",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == '-':
		return Token{
			Type:    TokenMinus,
			Lexeme:  "-",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == '+':
		return Token{
			Type:    TokenPlus,
			Lexeme:  "+",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == '=':
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
	case c == '!':
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
	case c == '<':
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
	case c == '>':
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
	case c == ';':
		return Token{
			Type:    TokenSemiColon,
			Lexeme:  ";",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == '/':
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

	case c == '*':
		return Token{
			Type:    TokenStar,
			Lexeme:  "*",
			Literal: "null",
			Line:    l.Line,
		}, nil, nil
	case c == '\n':
		l.Line += 1
		return l.NextToken()
	case c == ' ':
		return l.NextToken()
	case c == '\t':
		return l.NextToken()
	case c == '\r':
		return l.NextToken()
	case c == '"':
		l.Lexeme.Reset()
		l.Lexeme.WriteByte('"')
		for {
			nxt, _, err := r.ReadRune()
			if err == io.EOF {
				// if we reach the end of the file before terminating the string with '"'
				us := UnterminatedStringError{
					Line: l.Line,
				}
				us.Message = us.FormatMessage()
				return Token{}, us, nil
			}
			if err != nil {
				return Token{}, nil, err
			}
			_, err = l.Lexeme.WriteRune(nxt)
			if err != nil {
				return Token{}, nil, err
			}
			if nxt == '\n' {
				l.Line++
			}
			if nxt == '"' {
				// we have reached the end of the string literal
				bufBytes := l.Lexeme.Bytes()
				lexeme := string(bufBytes)
				return Token{
					Type:    TokenStringLiteral,
					Lexeme:  lexeme,
					Literal: strings.Trim(lexeme, "\""),
					Line:    l.Line,
				}, nil, nil
			}
		}
	case unicode.IsDigit(c):
		hasDot := false
		l.Lexeme.Reset()
		_, err := l.Lexeme.WriteRune(c)
		if err != nil {
			return Token{}, nil, err
		}
		for {
			nxt, _, err := r.ReadRune()
			if err == io.EOF {
				strNum := l.Lexeme.String()
				parsedFloat, err := strconv.ParseFloat(strNum, 64)
				if err != nil {
					return Token{}, nil, err
				}
				literal := strconv.FormatFloat(parsedFloat, 'f', -1, 64)
				if !strings.Contains(literal, ".") {
					literal += ".0"
				}
				return Token{
					Type:    TokenNumberLiteral,
					Lexeme:  strNum,
					Literal: literal,
					Line:    l.Line,
				}, nil, nil
			}
			if err != nil {
				return Token{}, nil, err
			}
			if !(unicode.IsDigit(nxt) || nxt == '.' && !hasDot) {
				err = r.UnreadRune()
				if err != nil {
					return Token{}, nil, err
				}
				strNum := l.Lexeme.String()
				parsedFloat, err := strconv.ParseFloat(strNum, 64)
				if err != nil {
					return Token{}, nil, err
				}
				literal := strconv.FormatFloat(parsedFloat, 'f', -1, 64)
				if !strings.Contains(literal, ".") {
					literal += ".0"
				}
				return Token{
					Type:    TokenNumberLiteral,
					Lexeme:  strNum,
					Literal: literal,
					Line:    l.Line,
				}, nil, nil
			}
			_, err = l.Lexeme.WriteRune(nxt)
			if err != nil {
				return Token{}, nil, err
			}
			if nxt == '.' {
				hasDot = true
			}
		}
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
		toks = append(toks, tok)
	}
}
