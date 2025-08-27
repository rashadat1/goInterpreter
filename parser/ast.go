package parser

import (
	"fmt"
	"strings"

	"github/goInterpreter/lexer"
)

// generic Expr[T] with generic Accept method
// generic ExprVisitor[T] interface with VisitX
// for every concrete Expr[T] (nodes)

type ExprVisitor[T, V any] interface {
	VisitBinary(*Binary[T, V]) V
	VisitUnary(*Unary[T, V]) V
	VisitGrouping(*Grouping[T, V]) V
	VisitLiteral(*Literal[T, V]) V
	VisitComma(*Comma[T, V]) V
	VisitTernary(*Ternary[T, V]) V
}

type Expr[T, V any] interface {
	// T: the type of literal values
	// V: the return type of the visitor
	Accept(ExprVisitor[T, V]) V
}

type Binary[T, V any] struct {
	Left     Expr[T, V]
	Operator lexer.Token
	Right    Expr[T, V]
}

func (bin *Binary[T, V]) Accept(visitor ExprVisitor[T, V]) V {
	return visitor.VisitBinary(bin)
}

type Unary[T, V any] struct {
	Operator lexer.Token
	Right    Expr[T, V]
}

func (un *Unary[T, V]) Accept(visitor ExprVisitor[T, V]) V {
	return visitor.VisitUnary(un)
}

type Grouping[T, V any] struct {
	Expression Expr[T, V]
}

func (gr *Grouping[T, V]) Accept(visitor ExprVisitor[T, V]) V {
	return visitor.VisitGrouping(gr)
}

type Literal[T, V any] struct {
	Value T
	Type  string
}

func (l *Literal[T, V]) Accept(visitor ExprVisitor[T, V]) V {
	return visitor.VisitLiteral(l)
}

type Comma[T, V any] struct {
	Left  Expr[T, V]
	Right Expr[T, V]
}

func (c *Comma[T, V]) Accept(visitor ExprVisitor[T, V]) V {
	return visitor.VisitComma(c)
}

type Ternary[T, V any] struct {
	Left   Expr[T, V]
	Middle Expr[T, V]
	Right  Expr[T, V]
}

func (t *Ternary[T, V]) Accept(visitor ExprVisitor[T, V]) V {
	return visitor.VisitTernary(t)
}

type AstPrinter struct{}

// printer should return a string so it implementst the Expr[T=string] interface
func (astp AstPrinter) VisitBinary(bin *Binary[any, string]) string {
	return printHelper(astp, bin.Operator.Lexeme, bin.Left, bin.Right)
}

func (astp AstPrinter) VisitLiteral(lit *Literal[any, string]) string {
	if lit.Value == nil {
		return "nil"
	}
	return fmt.Sprint(lit.Value)
}
func (astp AstPrinter) VisitGrouping(gr *Grouping[any, string]) string {
	return printHelper(astp, "group", gr.Expression)
}
func (astp AstPrinter) VisitUnary(un *Unary[any, string]) string {
	return printHelper(astp, un.Operator.Lexeme, un.Right)
}
func (astp AstPrinter) VisitComma(c *Comma[any, string]) string {
	return printHelper(astp, ",", c.Left, c.Right)
}
func (astp AstPrinter) VisitTernary(t *Ternary[any, string]) string {
	return printHelper(astp, "?:", t.Left, t.Middle, t.Right)
}
func printHelper(astp ExprVisitor[any, string], operation string, exprArgs ...Expr[any, string]) string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(operation)
	for _, exp := range exprArgs {
		sb.WriteString(" ")
		sb.WriteString(exp.Accept(astp))
	}
	sb.WriteString(")")
	return sb.String()
}
