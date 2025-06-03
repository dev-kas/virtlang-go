package evaluator

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/debugger"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

func evalMemberExpr(node *ast.MemberExpr, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	obj, err := Evaluate(node.Object, env, dbgr)
	if err != nil {
		return nil, err
	}

	if obj.Type != shared.Object && obj.Type != shared.Array && obj.Type != shared.ClassInstance {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of non-object or non-array (attempting to access properties of %v).", shared.Stringify(obj.Type)),
		}
	}

	if obj.Type == shared.Object {
		return evalMemberExpr_object(node, env, obj, dbgr)
	} else if obj.Type == shared.Array {
		return evalMemberExpr_array(node, env, dbgr)
	} else {
		return evalMemberExpr_class(node, env, dbgr)
	}
}

func evalMemberExpr_object(node *ast.MemberExpr, env *environment.Environment, obj *shared.RuntimeValue, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	var prop *shared.RuntimeValue
	if node.Computed {
		val, err := Evaluate(node.Value, env, dbgr)
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
		nilValue := values.MK_NIL()
		return &nilValue, nil
	}

	return obj.Value.(map[string]*shared.RuntimeValue)[key], nil
}

func evalMemberExpr_array(node *ast.MemberExpr, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	if !node.Computed {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of array by non-number (attempting to access properties by %v).", node.Value.GetType()),
		}
	}

	val, err := Evaluate(node.Value, env, dbgr)
	if err != nil {
		return nil, err
	}

	if val.Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of array by non-number (attempting to access properties by %v).", shared.Stringify(val.Type)),
		}
	}
	index := int(val.Value.(float64))

	var updatedArr *shared.RuntimeValue

	updatedArr, err = Evaluate(node.Object, env, dbgr)
	if err != nil {
		return nil, err
	}

	if updatedArr.Type != shared.Array {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of non-array (attempting to access properties of %v).", shared.Stringify(updatedArr.Type)),
		}
	}

	if index < 0 || index >= len(updatedArr.Value.([]shared.RuntimeValue)) {
		nilValue := values.MK_NIL()
		return &nilValue, nil
	}

	result := &updatedArr.Value.([]shared.RuntimeValue)[index]
	return result, nil
}

func evalMemberExpr_class(node *ast.MemberExpr, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	obj, err := Evaluate(node.Object, env, dbgr)
	if err != nil {
		return nil, err
	}

	if obj.Type != shared.ClassInstance {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot access property of non-class instance (attempting to access properties of %v).", shared.Stringify(obj.Type)),
		}
	}

	instance := obj.Value.(values.ClassInstanceValue)
	var key string

	if node.Computed {
		val, err := Evaluate(node.Value, env, dbgr)
		if err != nil {
			return nil, err
		}
		if val.Type != shared.String {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Cannot access property of class instance by non-string (attempting to access properties by %v).", shared.Stringify(val.Type)),
			}
		}
		key = val.Value.(string)
	} else {
		key = node.Value.(*ast.Identifier).Symbol
	}

	// Access control
	if !instance.Publics[key] {
		nilValue := values.MK_NIL()
		return &nilValue, nil
	}

	// Lookup
	value, err := instance.Data.LookupVar(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		nilValue := values.MK_NIL()
		return &nilValue, nil
	}
	return value, nil
}
