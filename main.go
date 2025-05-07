package main

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/evaluator"
	"github.com/dev-kas/virtlang-go/v2/parser"
)

func main() {
	// Create a new parser
	p := parser.New()

	// Create a new environment (for storing variables)
	env := environment.NewEnvironment(nil)

	// Your code to execute
	code := "let x = 1.1\nlet y = 2.2\nx + y"

	// Parse the code to produce an AST
	program, synErr := p.ProduceAST(code)
	if synErr != nil {
		fmt.Printf("%v\n", synErr)
		return
	}

	// Evaluate the program
	result, runErr := evaluator.Evaluate(program, &env)
	if runErr != nil {
		fmt.Printf("Runtime error: %v\n", runErr)
		return
	}

	// Print the result
	fmt.Printf("Result: %v\n", result.Value)
}
