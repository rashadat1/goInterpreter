package runner

import (
	"fmt"
	"github/goInterpreter/lexer"
	"os"
)

type RunnerResult struct {
	Tokens   string
	ErrorTok string
	ExitCode int
}

func Run(lex *lexer.Lexer) RunnerResult {
	res := ""
	resErr := ""
	var exitCode int
	toks, tokErrs, err := lex.ScanTokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred processing tokens: %s", err.Error())
	}
	res = toks.ToString()
	resErr = tokErrs.ToString()
	if len(resErr) != 0 {
		exitCode = 65
	}
	return RunnerResult{
		Tokens:   res,
		ErrorTok: resErr,
		ExitCode: exitCode,
	}
}
