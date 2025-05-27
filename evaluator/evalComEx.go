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

// compareValues handles the actual comparison of two values based on their types
func compareValues(lhs, rhs *shared.RuntimeValue) (int, *errors.RuntimeError) {
	// Different types are not comparable with <, >, etc.
	if lhs.Type != rhs.Type {
		return 0, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot compare different types: %v and %v", lhs.Type, rhs.Type),
		}
	}

	// Same type comparisons
	switch lhs.Type {
	case shared.Number:
		lhsVal := lhs.Value.(float64)
		rhsVal := rhs.Value.(float64)
		if lhsVal < rhsVal {
			return -1, nil
		} else if lhsVal > rhsVal {
			return 1, nil
		}
		return 0, nil

	case shared.String:
		lhsVal := lhs.Value.(string)
		rhsVal := rhs.Value.(string)
		if lhsVal < rhsVal {
			return -1, nil
		} else if lhsVal > rhsVal {
			return 1, nil
		}
		return 0, nil

	case shared.Boolean:
		lhsVal := lhs.Value.(bool)
		rhsVal := rhs.Value.(bool)
		if !lhsVal && rhsVal { // false < true
			return -1, nil
		} else if lhsVal && !rhsVal {
			return 1, nil
		}
		return 0, nil

	default:
		// For other types, comparison is not supported
		return 0, &errors.RuntimeError{
			Message: fmt.Sprintf("Comparison not supported for type: %v", lhs.Type),
		}
	}
}



// compareEqual handles == and !=
func compareEqual(lhs, rhs *shared.RuntimeValue, negate bool) (*shared.RuntimeValue, *errors.RuntimeError) {
	// If both sides are the same reference, they are equal (unless negated)
	if lhs == rhs {
		res := values.MK_BOOL(!negate)
		return &res, nil
	}

	// If types are different, they can't be equal
	if lhs.Type != rhs.Type {
		res := values.MK_BOOL(negate)
		return &res, nil
	}

	switch lhs.Type {
	case shared.Nil:
		res := values.MK_BOOL(!negate)
		return &res, nil

	case shared.Number, shared.String, shared.Boolean:
		// For primitive types, compare values directly
		result := lhs.Value == rhs.Value
		if negate {
			result = !result
		}
		res := values.MK_BOOL(result)
		return &res, nil

	case shared.Function:
		// For functions, we compare the actual function values, not just the pointers
		// This ensures that when a function is assigned to a variable, the comparison works
		lhsFn, lhsOk := lhs.Value.(*values.FunctionValue)
		rhsFn, rhsOk := rhs.Value.(*values.FunctionValue)
		
		if !lhsOk || !rhsOk {
			res := values.MK_BOOL(negate)
			return &res, nil
		}

		// Compare the function pointers for equality
		result := lhsFn == rhsFn
		if negate {
			result = !result
		}
		res := values.MK_BOOL(result)
		return &res, nil

	case shared.Array, shared.Object, shared.Class, shared.ClassInstance, shared.NativeFN:
		// For other reference types, they are only equal if they are the same reference
		res := values.MK_BOOL(negate)
		return &res, nil

	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot compare values of type %v", lhs.Type),
		}
	}
}

// compareLess handles < and <= operators
func compareLess(lhs, rhs *shared.RuntimeValue, orEqual bool) (*shared.RuntimeValue, *errors.RuntimeError) {
	// Different types are not comparable with <, >, etc.
	if lhs.Type != rhs.Type {
		res := values.MK_BOOL(false)
		return &res, nil
	}

	cmp, err := compareValues(lhs, rhs)
	if err != nil {
		res := values.MK_BOOL(false)
		return &res, nil
	}

	result := cmp < 0 || (orEqual && cmp == 0)
	res := values.MK_BOOL(result)
	return &res, nil
}

// compareGreater handles > and >= operators
func compareGreater(lhs, rhs *shared.RuntimeValue, orEqual bool) (*shared.RuntimeValue, *errors.RuntimeError) {
	// Different types are not comparable with <, >, etc.
	if lhs.Type != rhs.Type {
		res := values.MK_BOOL(false)
		return &res, nil
	}

	cmp, err := compareValues(lhs, rhs)
	if err != nil {
		res := values.MK_BOOL(false)
		return &res, nil
	}

	result := cmp > 0 || (orEqual && cmp == 0)
	res := values.MK_BOOL(result)
	return &res, nil
}

func evalComEx(expression *ast.CompareExpr, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	lhs, err := Evaluate(expression.LHS, env, dbgr)
	if err != nil {
		return nil, err
	}

	rhs, err := Evaluate(expression.RHS, env, dbgr)
	if err != nil {
		return nil, err
	}

	switch expression.Operator {
	case ast.Equal:
		return compareEqual(lhs, rhs, false)
	case ast.NotEqual:
		return compareEqual(lhs, rhs, true)
	case ast.LessThan:
		return compareLess(lhs, rhs, false)
	case ast.LessThanEqual:
		return compareLess(lhs, rhs, true)
	case ast.GreaterThan:
		return compareGreater(lhs, rhs, false)
	case ast.GreaterThanEqual:
		return compareGreater(lhs, rhs, true)
	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Unknown comparison operator: %v.", expression.Operator),
		}
	}
}
