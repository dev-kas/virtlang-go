package evaluator

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalClassProperty(node *ast.ClassProperty, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	propertyName := node.Name
	propertyValue := node.Value

	var value shared.RuntimeValue = values.MK_NIL()

	if propertyValue != nil {
		val, err := Evaluate(propertyValue, env, dbgr)
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
