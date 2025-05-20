package evaluator

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

func evalArrayExpr(expr *ast.ArrayLiteral, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	elements := expr.Elements
	results := make([]shared.RuntimeValue, len(elements))

	for index, element := range elements {
		result, err := Evaluate(element, env)
		if err != nil {
			return nil, err
		}

		results[index] = *result
	}

	result := values.MK_ARRAY(results)
	return &result, nil
}
