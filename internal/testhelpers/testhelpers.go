package testhelpers

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/evaluator"
	"github.com/dev-kas/virtlang-go/v4/parser"
	"github.com/dev-kas/virtlang-go/v4/shared"
)

// MustParse parses source code and returns AST program or fails the test
func MustParse(t *testing.T, src string) *ast.Program {
	t.Helper()
	p := parser.New("test")
	prog, err := p.ProduceAST(src)
	if err != nil {
		t.Fatal(err)
	}
	return prog
}

// MustEval evaluates source code fully (parse -> eval) and returns the value
func MustEval(t *testing.T, src string) *shared.RuntimeValue {
	t.Helper()
	prog := MustParse(t, src)
	env := environment.NewEnvironment(nil)
	val, err := evaluator.Evaluate(prog, env, nil)
	if err != nil {
		t.Fatal(err)
	}
	return val
}

// EvalNode evaluates a manually constructed AST node directly
func EvalNode(t *testing.T, node ast.Expr) *shared.RuntimeValue {
	t.Helper()
	env := environment.NewEnvironment(nil)
	val, err := evaluator.Evaluate(node, env, nil)
	if err != nil {
		t.Fatal(err)
	}
	return val
}

// ExpectParseError checks that parsing fails for invalid sources
func ExpectParseError(t *testing.T, src string) {
	t.Helper()
	p := parser.New("test")
	_, err := p.ProduceAST(src)
	if err == nil {
		t.Fatal("Expected parse error")
	}
}
