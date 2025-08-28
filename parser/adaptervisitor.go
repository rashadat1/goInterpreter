package parser

import (
	"fmt"
	"github/goInterpreter/parser/exprVisitors"
)

func TransformToStringAST(expr exprVisitors.Expr[any, interface{}]) exprVisitors.Expr[any, string] {
	switch e := expr.(type) {
	case *exprVisitors.Binary[any, interface{}]:
		return &exprVisitors.Binary[any, string]{
			Left:     TransformToStringAST(e.Left),
			Operator: e.Operator,
			Right:    TransformToStringAST(e.Right),
		}
	case *exprVisitors.Unary[any, interface{}]:
		return &exprVisitors.Unary[any, string]{
			Operator: e.Operator,
			Right:    TransformToStringAST(e.Right),
		}
	case *exprVisitors.Grouping[any, interface{}]:
		return &exprVisitors.Grouping[any, string]{
			Expression: TransformToStringAST(e.Expression),
		}
	case *exprVisitors.Literal[any, interface{}]:
		return &exprVisitors.Literal[any, string]{
			Value: fmt.Sprintf("%v", e.Value),
			Type:  e.Type,
		}
	case *exprVisitors.Comma[any, interface{}]:
		return &exprVisitors.Comma[any, string]{
			Left:  TransformToStringAST(e.Left),
			Right: TransformToStringAST(e.Right),
		}
	case *exprVisitors.Ternary[any, interface{}]:
		return &exprVisitors.Ternary[any, string]{
			Left:   TransformToStringAST(e.Left),
			Middle: TransformToStringAST(e.Middle),
			Right:  TransformToStringAST(e.Right),
		}
	}
	panic("unknown expr type")
}
