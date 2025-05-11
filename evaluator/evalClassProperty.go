package evaluator

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

func evalClassProperty(node *ast.ClassProperty, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	propertyName := node.Name
	propertyValue := node.Value

	var value shared.RuntimeValue = values.MK_NIL()

	if propertyValue != nil {
		val, err := Evaluate(propertyValue, env)
		if err != nil {
			return nil, err
		}
		value = *val
	}

	valPointer, err := env.DeclareVar(propertyName, value, false)
	if err != nil {
		return nil, err
	}

	return valPointer, nil
}
