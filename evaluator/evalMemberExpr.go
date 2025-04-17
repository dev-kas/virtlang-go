package evaluator

import (
	"VirtLang/ast"
	"VirtLang/environment"
	"VirtLang/errors"
	"VirtLang/shared"
	"VirtLang/values"
	"fmt"
)

func evalMemberExpr(node *ast.MemberExpr, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	obj, err := Evaluate(node.Object, env)
	if err != nil {
		return nil, err
	}

	if obj.Type != shared.Object {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of non-object (attempting to access properties of %v).", obj.Type),
		}
	}

	var prop *shared.RuntimeValue
	if node.Computed {
		val, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}

		prop = val
	} else {
		prop = &shared.RuntimeValue{
			Type:  shared.String,
			Value: node.Value.(*ast.Identifier).Symbol,
		}
	}

	var key string

	if prop.Type == shared.ValueType(ast.IdentifierNode) {
		key = prop.Value.(string)
	} else {
		// key = prop.Value.(shared.RuntimeValue).Value.(string)
		switch v := prop.Value.(type) {
		case string:
			key = v
		case int:
			key = fmt.Sprintf("%v", v)
		default:
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Invalid property key type: %T", prop.Value),
			}
		}

	}

	if obj.Value.(map[string]*shared.RuntimeValue)[key] == nil {
		// return nil, &errors.RuntimeError{
		// 	Message: fmt.Sprintf("Property '%s' does not exist.", key),
		// }

		nilValue := values.MK_NIL()
		return &nilValue, nil
	}

	return obj.Value.(map[string]*shared.RuntimeValue)[key], nil
}
