package evaluator

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/helpers"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalLogicEx(expression *ast.LogicalExpr, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	var lhs *shared.RuntimeValue
	var err *errors.RuntimeError

	if expression.Operator == ast.LogicalNOT {
		if expression.LHS != nil {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Logical NOT operator can only be used without LHS, got %v.", expression.LHS),
			}
		}

		rhs, err := Evaluate(expression.RHS, env, dbgr)
		if err != nil {
			return nil, err
		}

		res := values.MK_BOOL(!helpers.IsTruthy(rhs))
		return &res, nil
	}

	if expression.LHS != nil {
		lhs, err = Evaluate(*expression.LHS, env, dbgr)
		if err != nil {
			return nil, err
		}
	}

	switch expression.Operator {
	case ast.LogicalAND:
		if !helpers.IsTruthy(lhs) {
			return lhs, nil // short circuit: return if false
		}
		rhs, err := Evaluate(expression.RHS, env, dbgr)
		if err != nil {
			return nil, err
		}
		return rhs, nil

	case ast.LogicalOR:
		if helpers.IsTruthy(lhs) {
			return lhs, nil // short circuit: return if true
		}
		rhs, err := Evaluate(expression.RHS, env, dbgr)
		if err != nil {
			return nil, err
		}
		return rhs, nil

	case ast.LogicalNilCoalescing:
		if lhs == nil || (lhs.Type == shared.Nil && lhs.Value == nil) {
			rhs, err := Evaluate(expression.RHS, env, dbgr)
			if err != nil {
				return nil, err
			}
			return rhs, nil
		}
		return lhs, nil

	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Unknown logical operator: %v.", expression.Operator),
		}
	}
}
