package evaluator

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/debugger"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
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
