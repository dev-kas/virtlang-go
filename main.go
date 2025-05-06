package main

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/evaluator"
	"github.com/dev-kas/VirtLang-Go/parser"
)

func main() {
	// Create a new parser
	p := parser.New()

	// Create a new environment (for storing variables)
	env := environment.NewEnvironment(nil)

	// Your code to execute
	code := `let b = 1.2.2`

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
