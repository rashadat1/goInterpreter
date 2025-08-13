package lexer

type TokenType int

const (
	TokenEOF TokenType = iota
	// braces, parentheses
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	// punctuation
	TokenComma
	TokenDot
	TokenSemiColon
	// mathematical operators
	TokenMinus
	TokenPlus
	TokenEqual
	TokenBangEqual
	TokenEqualEqual
	TokenLess
	TokenLessEqual
	TokenGreater
	TokenGreaterEqual
	// misc
	TokenSlash
	TokenStar
	TokenBang
	TokenIdentifier
	// literals
	TokenStringLiteral
	TokenNumberLiteral
	// keywords
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFor
	TokenFun
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile
)

func (t TokenType) ToString() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenLeftParen:
		return "LEFT_PAREN"
	case TokenRightParen:
		return "RIGHT_PAREN"
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
	case TokenBang:
		return "BANG"
	case TokenBangEqual:
		return "BANG_EQUAL"
	case TokenLess:
		return "LESS"
	case TokenLessEqual:
		return "LESS_EQUAL"
	case TokenGreater:
		return "GREATER"
	case TokenGreaterEqual:
		return "GREATER_EQUAL"
	case TokenStringLiteral:
		return "STRING"
	case TokenNumberLiteral:
		return "NUMBER"
	case TokenAnd:
		return "AND"
	case TokenClass:
		return "CLASS"
	case TokenElse:
		return "ELSE"
	case TokenFalse:
		return "FALSE"
	case TokenFor:
		return "FOR"
	case TokenFun:
		return "FUN"
	case TokenIf:
		return "IF"
	case TokenNil:
		return "NIL"
	case TokenOr:
		return "OR"
	case TokenPrint:
		return "PRINT"
	case TokenReturn:
		return "RETURN"
	case TokenSuper:
		return "SUPER"
	case TokenThis:
		return "THIS"
	case TokenTrue:
		return "TRUE"
	case TokenVar:
		return "VAR"
	case TokenWhile:
		return "WHILE"
	case TokenIdentifier:
		return "IDENTIFIER"
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
	if t.Type.ToString() == "" {
		return ""
	}
	return t.Type.ToString() + " " + t.Lexeme + " " + LiteralToString(t.Literal) + string('\n')
}

type TokenizedText []Token

func (tt TokenizedText) ToString() string {
	out := ""
	for _, tok := range tt {
		out += tok.TokToString()
	}
	return out
}
