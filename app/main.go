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
		fmt.Fprintf(os.Stderr, "Usage: ./your_program.sh tokenize\n")
		fmt.Fprintf(os.Stderr, "Usage: ./your_program.sh parse <filename>\n")
		fmt.Fprintf(os.Stderr, "Usage: ./your_program.sh parse\n")
		os.Exit(1)
	}
	command := os.Args[1]
	if len(os.Args) == 2 {
		if command == "tokenize" {
			replTokenize()
		} else if command == "parse" {
			replParse()
		}
		os.Exit(0)
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
		lexRes := runner.RunLexer(&lex)
		switch command {
		case "tokenize":
			lexRes.Print()
			os.Exit(lexRes.ExitCode)
		case "parse":
			parseRes := runner.RunParser(lexRes.Tokens)
			if parseRes.ExitCode != 0 {
				os.Exit(parseRes.ExitCode)
			}
			parseRes.Print()
			os.Exit(parseRes.ExitCode)
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
			os.Exit(1)
		}
	} else {
		fmt.Println("EOF  null")
	}

}

func replTokenize() {
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
		lexRes := runner.RunLexer(&lex)
		lexRes.Print()
	}
}

func replParse() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprintf(os.Stdout, "> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error occurred reading line: %s", err.Error())
				continue
			}
			// if scanner.Scan() == false but no error occured it means the user hit control + D
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
		runRes := runner.RunLexer(&lex)
		parseRes := runner.RunParser(runRes.Tokens)
		if parseRes.ExitCode == 0 {
			parseRes.Print()
		}
	}

}
