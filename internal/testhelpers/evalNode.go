package testhelpers

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/evaluator"
	"github.com/dev-kas/virtlang-go/v4/shared"
)

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
