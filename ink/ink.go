package ink

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hanzala211/goinkinterpreter/parser"
	"github.com/hanzala211/goinkinterpreter/scanner"
	"github.com/hanzala211/goinkinterpreter/token"
)

type Ink struct {
	hadError bool
}

func (i *Ink) RunFile(file string) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	i.run(bytes)
	if i.hadError {
		os.Exit(65)
	}
	return nil
}

func (i *Ink) run(bytes []byte) {
	scanner := scanner.NewScanner(string(bytes), i)
	scanner.ScanTokens()
	p := parser.NewParser(scanner.Tokens, i)
	expr, err := p.Parse()
	if err != nil {
		i.ReportParserError(err.(parser.ParserError))
		return
	}
	fmt.Println(expr.String())
}

func (i *Ink) ReportError(line int, err error) {
	i.report(line, "", err)
}

func (i *Ink) report(line int, where string, err error) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	i.hadError = true
}

func (i *Ink) RunPrompt() {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if ok := input.Scan(); !ok {
			break
		}
		line := input.Text()
		i.run([]byte(line))
		i.hadError = false
	}
}

func (i *Ink) ReportParserError(err parser.ParserError) {
	if err.Token.Type == token.TokenType_EOF {
		i.report(err.Token.Line, " at end", err)
	} else {
		i.report(err.Token.Line, " at '"+err.Token.Lexeme+"'", err)
	}
}
