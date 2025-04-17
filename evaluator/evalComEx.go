package evaluator

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

func evalComEx(expression *ast.CompareExpr, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	lhs, err := Evaluate(expression.LHS, env)
	if err != nil {
		return nil, err
	}

	rhs, err := Evaluate(expression.RHS, env)
	if err != nil {
		return nil, err
	}

	if lhs.Type != shared.Number || rhs.Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot compare non-number types (%v and %v).", lhs.Type, rhs.Type),
		}
	}

	switch expression.Operator {
	case ast.Equal:
		result := values.MK_BOOL(lhs.Value.(int) == rhs.Value.(int))
		return &result, nil
	case ast.NotEqual:
		result := values.MK_BOOL(lhs.Value.(int) != rhs.Value.(int))
		return &result, nil
	case ast.LessThan:
		result := values.MK_BOOL(lhs.Value.(int) < rhs.Value.(int))
		return &result, nil
	case ast.LessThanEqual:
		result := values.MK_BOOL(lhs.Value.(int) <= rhs.Value.(int))
		return &result, nil
	case ast.GreaterThan:
		result := values.MK_BOOL(lhs.Value.(int) > rhs.Value.(int))
		return &result, nil
	case ast.GreaterThanEqual:
		result := values.MK_BOOL(lhs.Value.(int) >= rhs.Value.(int))
		return &result, nil
	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Unknown comparison operator: %v.", expression.Operator),
		}
	}
}
