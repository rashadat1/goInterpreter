package parser

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"

	"github/goInterpreter/lexer"
)

type Interpreter struct {
	HadError bool
}

func (i *Interpreter) VisitBinary(b *Binary[any, interface{}]) interface{} {
	left := i.evaluate(b.Left)
	if err, ok := left.(error); ok {
		return err
	}
	right := i.evaluate(b.Right)
	if err, ok := right.(error); ok {
		return err
	}
	operator := b.Operator
	// operators that can have non-numeric operands
	if operator.Type == lexer.TokenPlus {
		if checkStringOperands(left, right) {
			return left.(string) + right.(string)
		}
	}
	if operator.Type == lexer.TokenEqualEqual {
		ok, err := isEqual(left, right)
		if err != nil {
			errStr := fmt.Sprintf("[line %d] Error: invalid operation %v %s %v (mismatched types %s and %s)", operator.Line, left, operator.Lexeme, right, reflect.TypeOf(left).String(), reflect.TypeOf(right).String())
			return errors.New(errStr)
		}
		return ok
	}
	if operator.Type == lexer.TokenBangEqual {
		ok, err := isEqual(left, right)
		if err != nil {
			errStr := fmt.Sprintf("[line %d] Error: invalid operation %v %s %v (mismatched types %s and %s)", operator.Line, left, operator.Lexeme, right, reflect.TypeOf(left).String(), reflect.TypeOf(right).String())
			return errors.New(errStr)
		}
		return !ok
	}
	// all operators below are only defined for numeric operands (except TokenPlus which we already checked the string case)
	if !checkNumberOperands(left, right) {
		errStr := fmt.Sprintf("[line %d] Error: invalid operation %v %s %v (mismatched types %s and %s)", operator.Line, left, operator.Lexeme, right, reflect.TypeOf(left).String(), reflect.TypeOf(right).String())
		return errors.New(errStr)
	}
	switch operator.Type {
	case lexer.TokenPlus:
		return left.(float64) + right.(float64)
	case lexer.TokenMinus:
		return left.(float64) - right.(float64)

	case lexer.TokenSlash:
		return left.(float64) / right.(float64)

	case lexer.TokenStar:
		return left.(float64) * right.(float64)

	case lexer.TokenStarStar:
		return math.Pow(left.(float64), right.(float64))

	case lexer.TokenGreater:
		return left.(float64) > right.(float64)

	case lexer.TokenGreaterEqual:
		return left.(float64) >= right.(float64)

	case lexer.TokenLess:
		return left.(float64) < right.(float64)

	case lexer.TokenLessEqual:
		return left.(float64) <= right.(float64)
	}
	return nil
}
func (i *Interpreter) VisitUnary(un *Unary[any, interface{}]) interface{} {
	operator := un.Operator
	operand := i.evaluate(un.Right)
	if err, ok := operand.(error); ok {
		return err
	}
	switch operator.Type {
	case lexer.TokenMinus:
		if val, ok := operand.(float64); ok {
			return -val
		}
		errStr := fmt.Sprintf("[line %d] Error: invalid operation %v %s (operand must be numeric cannot be %s)", operator.Line, operand, operator.Lexeme, reflect.TypeOf(operand).String())
		return errors.New(errStr)
	case lexer.TokenBang:
		return !isTruthy(operand)
	}
	return nil
}
func (i *Interpreter) VisitGrouping(gr *Grouping[any, interface{}]) interface{} {
	groupExp := i.evaluate(gr.Expression)
	if err, ok := groupExp.(error); ok {
		return err
	}
	return groupExp
}
func (i *Interpreter) VisitLiteral(lit *Literal[any, interface{}]) interface{} {
	switch {
	case lit.Type == "number":
		strVal := lit.Value.(string)
		parsedFloatNum, _ := strconv.ParseFloat(strVal, 64)
		return parsedFloatNum
	case lit.Type == "string":
		strVal := lit.Value.(string)
		return strVal
	default:
		return lit.Value
	}
}
func (i *Interpreter) VisitComma(c *Comma[any, interface{}]) interface{} {
	left := i.evaluate(c.Left)
	if err, ok := left.(error); ok {
		return err
	}
	right := i.evaluate(c.Right)
	if err, ok := right.(error); ok {
		return err
	}
	return right
}
func (i *Interpreter) VisitTernary(t *Ternary[any, interface{}]) interface{} {
	condition := i.evaluate(t.Left)
	if err, ok := condition.(error); ok {
		return err
	}
	if isTruthy(condition) {
		ifTrue := i.evaluate(t.Middle)
		if err, ok := ifTrue.(error); ok {
			return err
		}
		return ifTrue
	}
	ifFalse := i.evaluate(t.Right)
	if err, ok := ifFalse.(error); ok {
		return err
	}
	return ifFalse
}
func (i *Interpreter) evaluate(expr Expr[any, interface{}]) interface{} {
	return expr.Accept(i)
}
func (i *Interpreter) Interpret(expr Expr[any, interface{}]) interface{} {
	v := i.evaluate(expr)
	if err, ok := v.(error); ok {
		log.Print(err.Error())
		i.HadError = true
		return err
	} else {
		return stringify(v)
	}
}
func stringify(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case nil:
		return "nil"
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}
func checkNumberOperands(op1, op2 interface{}) bool {
	_, ok1 := op1.(float64)
	_, ok2 := op2.(float64)
	return ok1 && ok2
}

func checkStringOperands(op1, op2 interface{}) bool {
	_, ok1 := op1.(string)
	_, ok2 := op2.(string)
	return ok1 && ok2
}
func checkBooleanOperands(op1, op2 interface{}) bool {
	_, ok1 := op1.(bool)
	_, ok2 := op2.(bool)
	return ok1 && ok2
}
func isEqual(op1, op2 interface{}) (bool, error) {
	switch {
	case op1 == nil && op2 == nil:
		return true, nil
	case op1 == nil || op2 == nil:
		return false, nil
	case checkBooleanOperands(op1, op2):
		return op1.(bool) == op2.(bool), nil
	case checkNumberOperands(op1, op2):
		return op1.(float64) == op2.(float64), nil
	case checkStringOperands(op1, op2):
		return op1.(string) == op2.(string), nil
	default:
		return false, errors.New("type mismatch")
	}
}
func isTruthy(op interface{}) bool {
	// nil and false are falsy and all other values are truthy
	if op == nil {
		return false
	}
	boolVal, ok := op.(bool)
	if ok {
		return boolVal
	}
	return true
}
