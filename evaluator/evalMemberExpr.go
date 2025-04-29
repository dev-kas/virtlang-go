package evaluator

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

func evalMemberExpr(node *ast.MemberExpr, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	obj, err := Evaluate(node.Object, env)
	if err != nil {
		return nil, err
	}

	if obj.Type != shared.Object && obj.Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of non-object or non-array (attempting to access properties of %v).", shared.Stringify(obj.Type)),
		}
	}

	if obj.Type == shared.Object {
		return evalMemberExpr_object(node, env, obj)
	} else {
		return evalMemberExpr_array(node, env, obj)
	}
}

func evalMemberExpr_object(node *ast.MemberExpr, env *environment.Environment, obj *shared.RuntimeValue) (*shared.RuntimeValue, *errors.RuntimeError) {
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
			// key = v[1 : len(v)-1]
		case int:
			key = fmt.Sprintf("%v", v)
		default:
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Invalid property key type: %T", prop.Value),
			}
		}

	}

	if obj.Value.(map[string]*shared.RuntimeValue)[key] == nil {
		nilValue := values.MK_NIL()
		return &nilValue, nil
	}

	return obj.Value.(map[string]*shared.RuntimeValue)[key], nil
}

func evalMemberExpr_array(node *ast.MemberExpr, env *environment.Environment, arr *shared.RuntimeValue) (*shared.RuntimeValue, *errors.RuntimeError) {
	if !node.Computed {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of array by non-number (attempting to access properties by %v).", node.Value.GetType()),
		}
	}

	val, err := Evaluate(node.Value, env)
	if err != nil {
		return nil, err
	}

	if val.Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of array by non-number (attempting to access properties by %v).", shared.Stringify(val.Type)),
		}
	}
	index := val.Value.(int)

	updatedArr, err := env.LookupVar(node.Object.(*ast.Identifier).Symbol)
	if err != nil {
		return nil, err
	}

	if index < 0 || index >= len(updatedArr.Value.([]shared.RuntimeValue)) {
		nilValue := values.MK_NIL()
		return &nilValue, nil
	}

	result := &updatedArr.Value.([]shared.RuntimeValue)[index]
	return result, nil
}
