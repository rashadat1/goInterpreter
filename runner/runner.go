package runner

import (
	"fmt"
	"github/goInterpreter/lexer"
	"github/goInterpreter/parser"
	"os"
)

type LexerResult struct {
	Tokens   lexer.TokenizedText
	ErrorTok lexer.TokenErrors
	ExitCode int
}
type ParserResult struct {
	Expr     parser.Expr[any, string]
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
	astp := parser.AstPrinter{}

	fmt.Println(expr.Accept(astp))
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
