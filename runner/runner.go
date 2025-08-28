package runner

import (
	"fmt"
	"os"

	"github/goInterpreter/lexer"
	"github/goInterpreter/parser"
	"github/goInterpreter/parser/exprVisitors"
)

type LexerResult struct {
	Tokens   lexer.TokenizedText
	ErrorTok lexer.TokenErrors
	ExitCode int
}
type ParserResult struct {
	Expr     exprVisitors.Expr[any, interface{}]
	Error    error
	ExitCode int
}
type EvaluateResult struct {
	Result   string
	Error    error
	ExitCode int
}

func (lr *LexerResult) Print() {
	res := ""
	resErr := ""
	res = lr.Tokens.ToString()
	resErr = lr.ErrorTok.ToString()
	fmt.Print(res)
	fmt.Fprint(os.Stderr, resErr)
}
func (pr *ParserResult) Print() {
	expr := pr.Expr
	astp := exprVisitors.AstPrinter{}
	astPrintInput := parser.TransformToStringAST(expr)
	fmt.Println(astPrintInput.Accept(astp))
}
func (er *EvaluateResult) Print() {
	if er.Error != nil {
		fmt.Fprint(os.Stderr, er.Error.Error())
	} else {
		fmt.Print(er.Result + string('\n'))
	}
}
func RunLexer(lex *lexer.Lexer) LexerResult {
	exitCode := 0
	toks, tokErrs, err := lex.ScanTokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred processing tokens: %s", err.Error())
	}
	if len(tokErrs) != 0 {
		exitCode = 65
	}
	return LexerResult{
		Tokens:   toks,
		ErrorTok: tokErrs,
		ExitCode: exitCode,
	}
}
func RunParser(tokens lexer.TokenizedText) ParserResult {
	exitCode := 0
	parser := parser.Parser{
		Tokens:   tokens,
		Position: 0,
	}
	expr := parser.Parse()
	if parser.HadError {
		exitCode = 65
	}
	return ParserResult{
		Expr:     expr,
		ExitCode: exitCode,
	}
}
func RunEvaluator(expr exprVisitors.Expr[any, interface{}]) EvaluateResult {
	evaluator := exprVisitors.Interpreter{
		HadError: false,
	}
	evalRes := EvaluateResult{}
	value := evaluator.Interpret(expr)
	switch res := value.(type) {
	case error:
		evalRes.Error = res
		evalRes.ExitCode = 70
	case string:
		evalRes.ExitCode = 0
		evalRes.Result = res
	}
	return evalRes
}
