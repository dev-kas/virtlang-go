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

func evalComEx(expression *ast.CompareExpr, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	lhs, err := Evaluate(expression.LHS, env, dbgr)
	if err != nil {
		return nil, err
	}

	rhs, err := Evaluate(expression.RHS, env, dbgr)
	if err != nil {
		return nil, err
	}

	if lhs.Type != shared.Number || rhs.Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot compare non-number types (%v and %v).", shared.Stringify(lhs.Type), shared.Stringify(rhs.Type)),
		}
	}

	switch expression.Operator {
	case ast.Equal:
		result := values.MK_BOOL(lhs.Value.(float64) == rhs.Value.(float64))
		return &result, nil
	case ast.NotEqual:
		result := values.MK_BOOL(lhs.Value.(float64) != rhs.Value.(float64))
		return &result, nil
	case ast.LessThan:
		result := values.MK_BOOL(lhs.Value.(float64) < rhs.Value.(float64))
		return &result, nil
	case ast.LessThanEqual:
		result := values.MK_BOOL(lhs.Value.(float64) <= rhs.Value.(float64))
		return &result, nil
	case ast.GreaterThan:
		result := values.MK_BOOL(lhs.Value.(float64) > rhs.Value.(float64))
		return &result, nil
	case ast.GreaterThanEqual:
		result := values.MK_BOOL(lhs.Value.(float64) >= rhs.Value.(float64))
		return &result, nil
	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Unknown comparison operator: %v.", expression.Operator),
		}
	}
}
