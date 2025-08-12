package parser

import (
	"errors"
	"log"

	"github/goInterpreter/lexer"
)

type Parser struct {
	Tokens   []lexer.Token
	Position int
}

func (p *Parser) Expression() Expr[any, any] {
	return p.Equality()
}
func (p *Parser) Equality() Expr[any, any] {
	newExpr := p.Comparison()
	for {
		if p.Match(lexer.TokenBangEqual, lexer.TokenEqualEqual) {
			// check if current token is of one of the
			// types in the args and advances
			operator := p.Previous()
			rightExpr := p.Comparison()
			newExpr = &Binary[any, any]{
				Left:     newExpr,
				Operator: operator,
				Right:    rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr
}
func (p *Parser) Comparison() Expr[any, any] {
	newExpr := p.Term()
	for {
		if p.Match(lexer.TokenGreater, lexer.TokenGreaterEqual, lexer.TokenLess, lexer.TokenLessEqual) {
			operator := p.Previous()
			rightExpr := p.Term()
			newExpr = &Binary[any, any]{
				Left:     newExpr,
				Operator: operator,
				Right:    rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr
}
func (p *Parser) Term() Expr[any, any] {
	newExpr := p.Factor()
	for {
		if p.Match(lexer.TokenMinus, lexer.TokenPlus) {
			operator := p.Previous()
			rightExpr := p.Factor()
			newExpr = &Binary[any, any]{
				Left:     newExpr,
				Operator: operator,
				Right:    rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr
}
func (p *Parser) Factor() Expr[any, any] {
	newExpr := p.Unary()
	for {
		if p.Match(lexer.TokenSlash, lexer.TokenStar) {
			operator := p.Previous()
			rightExpr := p.Unary()
			newExpr = &Binary[any, any]{
				Left:     newExpr,
				Operator: operator,
				Right:    rightExpr,
			}
		} else {
			break
		}
	}
	return newExpr
}
func (p *Parser) Unary() Expr[any, any] {
	if p.Match(lexer.TokenBang, lexer.TokenMinus) {
		operator := p.Previous()
		expr := p.Unary()
		return &Unary[any, any]{
			Operator: operator,
			Right:    expr,
		}
	}
	return p.Primary()
}
func (p *Parser) Primary() Expr[any, any] {
	if p.Match(lexer.TokenLeftParen) {
		expr := p.Expression()
		p.Consume(lexer.TokenRightParen, "Expect ')' after expression.")
		return &Grouping[any, any]{
			Expression: expr,
		}
	}
	if p.Match(lexer.TokenTrue) {
		return &Literal[any, any]{
			Value: true,
		}
	}
	if p.Match(lexer.TokenFalse) {
		return &Literal[any, any]{
			Value: false,
		}
	}
	if p.Match(lexer.TokenNil) {
		return &Literal[any, any]{
			Value: nil,
		}
	}
	if p.Match(lexer.TokenStringLiteral) {
		return &Literal[any, any]{
			Value: p.Previous().Literal,
		}
	}
	if p.Match(lexer.TokenNumberLiteral) {
		return &Literal[any, any]{
			Value: p.Previous().Literal,
		}
	}
	log.Fatal("Invalid expression")
	return nil
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
func (p *Parser) Match(tokenTypeArgs ...lexer.TokenType) bool {
	// check if current token's type matches any of the tokenTypes in the args
	if p.IsAtEnd() {
		return false
	}
	for _, tokenType := range tokenTypeArgs {
		if p.Peek().Type == tokenType {
			// consume token if we found a match
			p.Advance()
			return true
		}
	}
	return false
}
func (p *Parser) Advance() lexer.Token {
	if !p.IsAtEnd() {
		p.Position++
	}
	return p.Previous()
}
func (p *Parser) Consume(tknType lexer.TokenType, message string) lexer.Token {
	// look for closing ) if no closing ) then report an error
	if p.IsAtEnd() || p.Peek().Type != tknType {
		log.Fatal(errors.New(message))
	}
	return p.Advance()
}
