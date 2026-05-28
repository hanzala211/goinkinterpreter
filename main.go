package main

import (
	"fmt"
	"os"

	"github.com/hanzala211/goinkinterpreter/ink"
)

func main() {
	ink := ink.NewInk()
	if len(os.Args) > 2 {
		fmt.Println("Usage: go run main.go <file>")
		os.Exit(64)
	}
	if len(os.Args) == 2 {
		ink.RunFile(os.Args[1])
	} else {
		ink.RunPrompt()
	}
}
