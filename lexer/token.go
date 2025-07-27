package lexer

type TokenType int

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
	TokenEqual
	TokenEqualEqual
	TokenSemiColon
	TokenSlash
	TokenStar
)

func (t TokenType) ToString() string {
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
	case TokenEqual:
		return "EQUAL"
	case TokenEqualEqual:
		return "EQUAL_EQUAL"
	default:
		return ""
	}
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
	return t.Type.ToString() + " " + t.Lexeme + " " + LiteralToString(t.Literal) + string('\n')
}

type TokenizedText []Token

func (tt TokenizedText) ToString() string {
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
