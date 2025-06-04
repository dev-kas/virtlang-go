package evaluator

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalArrayExpr(expr *ast.ArrayLiteral, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	elements := expr.Elements
	results := make([]shared.RuntimeValue, len(elements))

	for index, element := range elements {
		result, err := Evaluate(element, env, dbgr)
		if err != nil {
			return nil, err
		}

		results[index] = *result
	}

	result := values.MK_ARRAY(results)
	return &result, nil
}
