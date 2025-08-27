package parser

import "fmt"

func TransformToStringAST(expr Expr[any, interface{}]) Expr[any, string] {
	switch e := expr.(type) {
	case *Binary[any, interface{}]:
		return &Binary[any, string]{
			Left:     TransformToStringAST(e.Left),
			Operator: e.Operator,
			Right:    TransformToStringAST(e.Right),
		}
	case *Unary[any, interface{}]:
		return &Unary[any, string]{
			Operator: e.Operator,
			Right:    TransformToStringAST(e.Right),
		}
	case *Grouping[any, interface{}]:
		return &Grouping[any, string]{
			Expression: TransformToStringAST(e.Expression),
		}
	case *Literal[any, interface{}]:
		return &Literal[any, string]{
			Value: fmt.Sprintf("%v", e.Value),
			Type:  e.Type,
		}
	case *Comma[any, interface{}]:
		return &Comma[any, string]{
			Left:  TransformToStringAST(e.Left),
			Right: TransformToStringAST(e.Right),
		}
	case *Ternary[any, interface{}]:
		return &Ternary[any, string]{
			Left:   TransformToStringAST(e.Left),
			Middle: TransformToStringAST(e.Middle),
			Right:  TransformToStringAST(e.Right),
		}
	}
	panic("unknown expr type")
}
