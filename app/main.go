package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github/goInterpreter/lexer"
	"github/goInterpreter/runner"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: ./your_program.sh tokenize <filename>\n")
		fmt.Fprintf(os.Stderr, "Usage: ./your_program.sh tokenize")
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		replEcho()
		os.Exit(0)
	}
	command := os.Args[1]
	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
	fileName := os.Args[2]
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %s", fileName, err)
		os.Exit(1)
	}
	if len(fileContents) > 0 {
		r := bytes.NewReader(fileContents)
		lex := lexer.Lexer{
			Reader: r,
			Line:   1,
			Lexeme: bytes.NewBuffer(nil),
		}
		runRes := runner.Run(&lex)
		fmt.Print(runRes.Tokens)
		fmt.Fprint(os.Stderr, runRes.ErrorTok)
		os.Exit(runRes.ExitCode)
	} else {
		fmt.Println("EOF  null")
	}
}

func replEcho() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprintf(os.Stdout, "> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error occurred reading line: %s", err.Error())
				continue
			}
			// if scanner.Scan() == false but not error occurred it means user hit control + D
			break
		}
		line := scanner.Text()
		if line == "" || line == string('\n') {
			continue
		}
		lineBytes := []byte(line)
		r := bytes.NewReader(lineBytes)
		lex := lexer.Lexer{
			Reader: r,
			Line:   1,
			Lexeme: bytes.NewBuffer(nil),
		}
		runRes := runner.Run(&lex)
		fmt.Print(runRes.Tokens)
		fmt.Fprint(os.Stderr, runRes.ErrorTok)
		fmt.Printf("Input code: %s\n", line)
	}
}
