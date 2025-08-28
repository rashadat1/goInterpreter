package statements

import "github/goInterpreter/parser/exprVisitors"

type StmtVisitor[T, V any] interface {
	VisitExprStmt(*ExprStmt[T, V]) V
	VisitPrintStmt(*PrintStmt[T, V]) V
}

type Stmt[T, V any] interface {
	Accept(StmtVisitor[T, V]) V
}
type ExprStmt[T, V any] struct {
	Expr exprVisitors.Expr[T, V]
}
type PrintStmt[T, V any] struct {
	Expr exprVisitors.Expr[T, V]
}

func (est *ExprStmt[T, V]) Accept(visitor StmtVisitor[T, V]) V {
	return visitor.VisitExprStmt(est)
}
func (pst *PrintStmt[T, V]) Accept(visitor StmtVisitor[T, V]) V {
	return visitor.VisitPrintStmt(pst)
}
