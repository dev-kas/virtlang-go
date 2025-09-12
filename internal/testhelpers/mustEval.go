package testhelpers

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/evaluator"
	"github.com/dev-kas/virtlang-go/v4/shared"
)

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
