package parser

import (
	"errors"
	"log"
	"slices"

	"github/goInterpreter/lexer"
)

type Parser struct {
	Tokens   []lexer.Token
	Position int
	HadError bool
}

func (p *Parser) Parse() Expr[any, string] {
	expr, err := p.Expression()
	if err != nil {
		p.HadError = true
	}
	return expr
}

func (p *Parser) Expression() (Expr[any, string], error) {
	expr, err := p.Comma()
	if err != nil {
		return nil, err
	}
	return expr, nil
}
func (p *Parser) Comma() (Expr[any, string], error) {
	if p.MissingLeftOperand([]lexer.TokenType{lexer.TokenComma}) {
		right, _ := p.Equality()
		return right, nil
	}
	newExpr, err := p.Equality()
	if err != nil {
		return nil, err
	}
	for {
		if p.Match([]lexer.TokenType{lexer.TokenComma}) {
			rightExpr, err := p.Equality()
			if err != nil {
				return nil, err
			}
			newExpr = &Comma[any, string]{
				Left:  newExpr,
				Right: rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr, nil
}

// <-- Ternary Operator goes here
func (p *Parser) Equality() (Expr[any, string], error) {
	equalityOperators := []lexer.TokenType{lexer.TokenBangEqual, lexer.TokenEqualEqual}

	if p.MissingLeftOperand(equalityOperators) {
		right, _ := p.Comparison()
		return right, nil
	}
	newExpr, err := p.Comparison()
	if err != nil {
		return nil, err
	}
	for {
		if p.Match(equalityOperators) {
			// check if current token is of one of the
			// types in the args and advances
			operator := p.Previous()
			rightExpr, err := p.Comparison()
			if err != nil {
				return nil, err
			}
			newExpr = &Binary[any, string]{
				Left:     newExpr,
				Operator: operator,
				Right:    rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr, nil
}
func (p *Parser) Comparison() (Expr[any, string], error) {
	compareOperators := []lexer.TokenType{lexer.TokenGreater, lexer.TokenGreaterEqual, lexer.TokenLess, lexer.TokenLessEqual}

	if p.MissingLeftOperand(compareOperators) {
		right, _ := p.Term()
		return right, nil
	}
	newExpr, err := p.Term()
	if err != nil {
		return nil, err
	}
	for {
		if p.Match(compareOperators) {
			operator := p.Previous()
			rightExpr, err := p.Term()
			if err != nil {
				return nil, err
			}
			newExpr = &Binary[any, string]{
				Left:     newExpr,
				Operator: operator,
				Right:    rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr, nil
}
func (p *Parser) Term() (Expr[any, string], error) {
	termOperators := []lexer.TokenType{lexer.TokenPlus, lexer.TokenMinus}

	if p.MissingLeftOperand([]lexer.TokenType{lexer.TokenPlus}) {
		right, _ := p.Factor()
		return right, nil
	}
	newExpr, err := p.Factor()
	if err != nil {
		return nil, err
	}
	for {
		if p.Match(termOperators) {
			operator := p.Previous()
			rightExpr, err := p.Factor()
			if err != nil {
				return nil, err
			}
			newExpr = &Binary[any, string]{
				Left:     newExpr,
				Operator: operator,
				Right:    rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr, nil
}
func (p *Parser) Factor() (Expr[any, string], error) {
	factorOperators := []lexer.TokenType{lexer.TokenSlash, lexer.TokenStar}

	if p.MissingLeftOperand(factorOperators) {
		right, _ := p.Unary()
		return right, nil
	}
	newExpr, err := p.Unary()
	if err != nil {
		return nil, err
	}
	for {
		if p.Match(factorOperators) {
			operator := p.Previous()
			rightExpr, err := p.Unary()
			if err != nil {
				return nil, err
			}
			newExpr = &Binary[any, string]{
				Left:     newExpr,
				Operator: operator,
				Right:    rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr, nil
}
func (p *Parser) Unary() (Expr[any, string], error) {
	unaryOperators := []lexer.TokenType{lexer.TokenBang, lexer.TokenMinus}

	if p.Match(unaryOperators) {
		operator := p.Previous()
		expr, err := p.Unary()
		if err != nil {
			return nil, err
		}
		return &Unary[any, string]{
			Operator: operator,
			Right:    expr,
		}, nil
	}
	return p.Primary()
}
func (p *Parser) Primary() (Expr[any, string], error) {
	if p.Match([]lexer.TokenType{lexer.TokenLeftParen}) {
		// after matching an open parentheses we parse the expression inside of it
		// and log an error if the expression is not followed by a closing parentheses
		expr, err := p.Expression()
		if err != nil {
			return nil, err
		}
		_, err = p.Consume(lexer.TokenRightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &Grouping[any, string]{
			Expression: expr,
		}, nil
	}
	if p.Match([]lexer.TokenType{lexer.TokenTrue}) {
		return &Literal[any, string]{
			Value: true,
		}, nil
	}
	if p.Match([]lexer.TokenType{lexer.TokenFalse}) {
		return &Literal[any, string]{
			Value: false,
		}, nil
	}
	if p.Match([]lexer.TokenType{lexer.TokenNil}) {
		return &Literal[any, string]{
			Value: nil,
		}, nil
	}
	if p.Match([]lexer.TokenType{lexer.TokenStringLiteral}) {
		return &Literal[any, string]{
			Value: p.Previous().Literal,
		}, nil
	}
	if p.Match([]lexer.TokenType{lexer.TokenNumberLiteral}) {
		return &Literal[any, string]{
			Value: p.Previous().Literal,
		}, nil
	}
	parseError := ParserError{
		Line:    p.Peek().Line,
		Message: "Expecting expression",
	}
	msg := parseError.Report(p.Peek())
	return nil, errors.New(msg)
}

// Utility Methods
func (p *Parser) Previous() lexer.Token {
	// retrieve the last emitted token
	return p.Tokens[p.Position-1]
}
func (p *Parser) Peek() lexer.Token {
	// get the next token to be consumed without consuming it
	return p.Tokens[p.Position]
}
func (p *Parser) IsAtEnd() bool {
	// recall every token slice ends with an EOF token
	return p.Peek().Type == lexer.TokenEOF
}
func (p *Parser) Match(tokenTypeArgs []lexer.TokenType) bool {
	// check if current token's type matches any of the tokenTypes in the args
	if p.IsAtEnd() {
		return false
	}
	if slices.Contains(tokenTypeArgs, p.Peek().Type) {
		// consume token if it appears in the set of match types
		p.Advance()
		return true
	}
	return false
}
func (p *Parser) Advance() lexer.Token {
	if !p.IsAtEnd() {
		p.Position++
	}
	return p.Previous()
}
func (p *Parser) Consume(tknType lexer.TokenType, message string) (lexer.Token, error) {
	// look for closing ) if no closing ) then report an error
	if p.IsAtEnd() || p.Peek().Type != tknType {
		parseError := ParserError{
			Line:    p.Peek().Line,
			Message: message,
		}
		msg := parseError.Report(p.Peek())
		p.HadError = true
		return lexer.Token{}, errors.New(msg)
	}
	return p.Advance(), nil
}
func (p *Parser) MissingLeftOperand(operators []lexer.TokenType) bool {
	if slices.Contains(operators, p.Peek().Type) {
		parseError := ParserError{
			Line:    p.Peek().Line,
			Message: "Missing left-hand operand",
		}
		op := p.Advance()
		errMessage := parseError.Report(op)
		log.Print(errMessage)
		p.HadError = true
		return true
	}
	return false
}
