package evaluator

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
)

func evalVarAssignment(node *ast.VarAssignmentExpr, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	if node.Assignee.GetType() == ast.IdentifierNode {
		varname := node.Assignee.(*ast.Identifier).Symbol

		value, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}

		return env.AssignVar(varname, *value)
	} else if node.Assignee.GetType() == ast.MemberExprNode {
		memberExpr := node.Assignee.(*ast.MemberExpr)
		obj, err := Evaluate(memberExpr.Object, env)
		if err != nil {
			return nil, err
		}

		if obj.Type != shared.Object {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Cannot access property of non-object (attempting to access properties of %v).", obj.Type),
			}
		}

		var prop *shared.RuntimeValue
		if memberExpr.Computed {
			val, err := Evaluate(memberExpr.Value, env)
			if err != nil {
				return nil, err
			}

			prop = val
		} else {
			prop = &shared.RuntimeValue{
				Type:  shared.String,
				Value: memberExpr.Value.(*ast.Identifier).Symbol,
			}
		}

		var key string

		if prop.Type == shared.ValueType(ast.IdentifierNode) {
			identifier, ok := prop.Value.(*ast.Identifier)
			if !ok {
				return nil, &errors.RuntimeError{
					Message: "Cannot access property of object by non-string key.",
				}
			}
			key = identifier.Symbol
		} else {
			// key = prop.Value.(shared.RuntimeValue).Value.(string)
			key = prop.Value.(string)
		}

		value, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}

		obj.Value.(map[string]*shared.RuntimeValue)[key] = value

		return value, nil
	} else {
		return nil, &errors.RuntimeError{
			Message: "Cannot access property of object by non-string key.",
		}
	}
}
