package main

import (
	"fmt"
	"os"

	"github.com/hanzala211/goinkinterpreter/ink"
)

func main() {
	ink := ink.NewInk()
	if len(os.Args) > 2 {
		fmt.Println("Usage: goinkinterpreter <file>")
		os.Exit(64)
	}
	if len(os.Args) == 2 {
		err := ink.RunFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	} else {
		ink.RunPrompt()
	}
}
