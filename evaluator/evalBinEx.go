package evaluator

import (
	"fmt"
	"math"

	"github.com/dev-kas/virtlang-go/ast"
	"github.com/dev-kas/virtlang-go/environment"
	"github.com/dev-kas/virtlang-go/errors"
	"github.com/dev-kas/virtlang-go/shared"
	"github.com/dev-kas/virtlang-go/values"
)

func countDecimalPlaces(f float64) int {
	if f == float64(int(f)) {
		return 0
	}

	f = math.Abs(f)
	maxDecimalPlaces := 15 // float64 precision limit

	for i := 0; i < maxDecimalPlaces; i++ {
		f *= 10
		if f == math.Floor(f) {
			return i + 1
		}
	}

	return maxDecimalPlaces
}

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
			Message: fmt.Sprintf("Cannot perform binary operation on non-number types (%v and %v).", shared.Stringify(lhs.Type), shared.Stringify(rhs.Type)),
		}
	}

	lhsValue := lhs.Value.(float64)
	rhsValue := rhs.Value.(float64)

	lhsIsFloat := lhsValue != float64(int(lhsValue))
	rhsIsFloat := rhsValue != float64(int(rhsValue))

	lhsFloatingDigits := 0
	rhsFloatingDigits := 0

	if lhsIsFloat {
		lhsFloatingDigits = countDecimalPlaces(lhsValue)
		lhsValue = lhsValue * math.Pow(10, float64(lhsFloatingDigits))
	}

	if rhsIsFloat {
		rhsFloatingDigits = countDecimalPlaces(rhsValue)
		rhsValue = rhsValue * math.Pow(10, float64(rhsFloatingDigits))
	}

	result := values.MK_NIL()

	switch binOp.Operator {
	case ast.Plus:
		result = values.MK_NUMBER(lhsValue + rhsValue)

	case ast.Minus:
		result = values.MK_NUMBER(lhsValue - rhsValue)

	case ast.Multiply:
		result = values.MK_NUMBER(lhsValue * rhsValue)

	case ast.Divide:
		result = values.MK_NUMBER(lhsValue / rhsValue)

	case ast.Modulo:
		// Formula applied:
		// mod(a, b) = a - floor(a/b) * b
		result = values.MK_NUMBER(lhsValue - math.Floor(lhsValue/rhsValue)*rhsValue)

	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Unknown binary operator: %v.", binOp.Operator),
		}
	}

	if lhsIsFloat || rhsIsFloat {
		maxDecimalPlaces := max(lhsFloatingDigits, rhsFloatingDigits)
		result.Value = result.Value.(float64) / math.Pow(10, float64(maxDecimalPlaces))
	}

	return &result, nil
}
