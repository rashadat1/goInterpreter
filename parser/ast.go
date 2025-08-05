package parser

import (
	"github/goInterpreter/lexer"
)

// generic Expr[T] with generic Accept method
// generic ExprVisitor[T] interface with VisitX
// for every concrete Expr[T] (nodes)

type ExprVisitor[T any] interface {
	VisitBinary(*Binary[T]) T
	VisitUnary(*Unary[T]) T
	VisitGrouping(*Grouping[T]) T
	VisitLiteral(*Literal[T]) T
}

type Expr[T any] interface {
	Accept(ExprVisitor[T]) T
}

type Binary[T any] struct {
	Left     Expr[T]
	Operator lexer.Token
	Right    Expr[T]
}

func (bin *Binary[T]) Accept(visitor ExprVisitor[T]) T {
	return visitor.VisitBinary(bin)
}

type Unary[T any] struct {
	Operator lexer.Token
	Right    Expr[T]
}

func (un *Unary[T]) Accept(visitor ExprVisitor[T]) T {
	return visitor.VisitUnary(un)
}

type Grouping[T any] struct {
	Expression Expr[T]
}

func (gr *Grouping[T]) Accept(visitor ExprVisitor[T]) T {
	return visitor.VisitGrouping(gr)
}

type Literal[T any] struct {
	Value T
}

func (l *Literal[T]) Accept(visitor ExprVisitor[T]) T {
	return visitor.VisitLiteral(l)
}
