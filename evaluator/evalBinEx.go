package evaluator

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

func evalBinEx(binOp *ast.BinaryExpr, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	lhs, err := Evaluate(binOp.LHS, env)
	if err != nil {
		return nil, err
	}

	rhs, err := Evaluate(binOp.RHS, env)
	if err != nil {
		return nil, err
	}

	if lhs.Type != shared.Number || rhs.Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot perform binary operation on non-number types (%v and %v).", lhs.Type, rhs.Type),
		}
	}

	switch binOp.Operator {
	case ast.Plus:
		result := values.MK_NUMBER(lhs.Value.(int) + rhs.Value.(int))
		return &result, nil

	case ast.Minus:
		result := values.MK_NUMBER(lhs.Value.(int) - rhs.Value.(int))
		return &result, nil

	case ast.Multiply:
		result := values.MK_NUMBER(lhs.Value.(int) * rhs.Value.(int))
		return &result, nil

	case ast.Divide:
		result := values.MK_NUMBER(lhs.Value.(int) / rhs.Value.(int))
		return &result, nil

	case ast.Modulo:
		result := values.MK_NUMBER(lhs.Value.(int) % rhs.Value.(int))
		return &result, nil

	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Unknown binary operator: %v.", binOp.Operator),
		}
	}
}
