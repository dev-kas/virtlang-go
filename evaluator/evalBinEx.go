package evaluator

import (
	"fmt"
	"math"

	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

func countDecimalPlaces(f float64) int {
	if f == float64(int(f)) {
		return 0
	}

	f = math.Abs(f)

	maxDecimalPlaces := 15

	epsilon := 1e-13

	for i := 0; i < maxDecimalPlaces; i++ {
		f *= 10
		if math.Abs(f-math.Floor(f)) < epsilon {
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

	opr := binOp.Operator

	var result *shared.RuntimeValue

	switch opr {
	case ast.Plus, ast.Minus:
		result, err = plusMinus(lhs, rhs, opr == ast.Plus)
	case ast.Multiply:
		result, err = multiply(lhs, rhs)
	case ast.Divide:
		result, err = divide(lhs, rhs)
	case ast.Modulo:
		result, err = modulo(lhs, rhs)
	}

	return result, err
}

func plusMinus(lhs, rhs *shared.RuntimeValue, isAddition bool) (*shared.RuntimeValue, *errors.RuntimeError) {
	if isAddition {
		if (lhs.Type != shared.Number && lhs.Type != shared.String) ||
			(rhs.Type != shared.Number && rhs.Type != shared.String) {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Addition binary operation can only be performed on numbers and strings. Attempted to perform addition on `%s` and `%s`", shared.Stringify(lhs.Type), shared.Stringify(rhs.Type)),
			}
		}
	} else {
		if lhs.Type != shared.Number || rhs.Type != shared.Number {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("This binary operation can only be performed on numbers. Attempted to perform subtraction on `%s` and `%s`", shared.Stringify(lhs.Type), shared.Stringify(rhs.Type)),
			}
		}
	}

	if lhs.Type != rhs.Type {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Type mismatch: LHS and RHS must be of the same type. Attempted to perform binary operation on distinct types `%s` and `%s`", shared.Stringify(lhs.Type), shared.Stringify(rhs.Type)),
		}
	}

	var result shared.RuntimeValue
	var err *errors.RuntimeError

	if lhs.Type == shared.String && rhs.Type == shared.String && isAddition {
		result = values.MK_STRING(lhs.Value.(string)[:len(lhs.Value.(string))-1] + rhs.Value.(string)[1:])
	} else {
		lhsValue := lhs.Value.(float64)
		rhsValue := rhs.Value.(float64)

		lhsDecimalPlaces := countDecimalPlaces(lhsValue)
		rhsDecimalPlaces := countDecimalPlaces(rhsValue)
		maxDecimalPlaces := max(lhsDecimalPlaces, rhsDecimalPlaces)

		scalingFactor := math.Pow10(maxDecimalPlaces)
		lhsValue = lhsValue * scalingFactor
		rhsValue = rhsValue * scalingFactor

		r := float64(0)

		if isAddition {
			r = lhsValue + rhsValue
		} else {
			r = lhsValue - rhsValue
		}

		r = r / scalingFactor

		result = values.MK_NUMBER(r)
	}

	return &result, err
}

func multiply(lhs, rhs *shared.RuntimeValue) (*shared.RuntimeValue, *errors.RuntimeError) {
	if lhs.Type != shared.Number || rhs.Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("This binary operation can only be performed on numbers. Attempted to perform multiplication on `%s` and `%s`", shared.Stringify(lhs.Type), shared.Stringify(rhs.Type)),
		}
	}

	if rhs.Value.(float64) == 0 || lhs.Value.(float64) == 0 {
		result := values.MK_NUMBER(0)
		return &result, nil
	}

	result := values.MK_NUMBER(0)

	lhsValue := lhs.Value.(float64)
	rhsValue := rhs.Value.(float64)

	lhsDecimalPlaces := countDecimalPlaces(lhsValue)
	rhsDecimalPlaces := countDecimalPlaces(rhsValue)

	lhsScalingFactor := math.Pow10(lhsDecimalPlaces)
	rhsScalingFactor := math.Pow10(rhsDecimalPlaces)

	lhsValue = lhsValue * lhsScalingFactor
	rhsValue = rhsValue * rhsScalingFactor

	r := float64(0)
	r = lhsValue * rhsValue
	r = r / (lhsScalingFactor * rhsScalingFactor)

	result = values.MK_NUMBER(r)

	return &result, nil
}

func modulo(lhs, rhs *shared.RuntimeValue) (*shared.RuntimeValue, *errors.RuntimeError) {
	if lhs.Type != shared.Number || rhs.Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("This binary operation can only be performed on numbers. Attempted to perform modulo on `%s` and `%s`", shared.Stringify(lhs.Type), shared.Stringify(rhs.Type)),
		}
	}

	if rhs.Value.(float64) == 0 {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot perform modulo operation with zero as the divisor. Attempted to perform modulo on `%g` with `0` as the divisor", lhs.Value.(float64)),
		}
	}

	lhsValue := lhs.Value.(float64)
	rhsValue := rhs.Value.(float64)

	lhsDecimalPlaces := countDecimalPlaces(lhsValue)
	rhsDecimalPlaces := countDecimalPlaces(rhsValue)
	maxDecimalPlaces := max(lhsDecimalPlaces, rhsDecimalPlaces)

	scalingFactor := math.Pow10(maxDecimalPlaces)
	lhsValue = lhsValue * scalingFactor
	rhsValue = rhsValue * scalingFactor

	r := float64(0)

	// Formula applied:
	// modulo(a, b) = a - floor(a/b) * b
	r = lhsValue - math.Floor(lhsValue/rhsValue)*rhsValue

	r = r / scalingFactor

	result := values.MK_NUMBER(r)

	return &result, nil
}

func divide(lhs, rhs *shared.RuntimeValue) (*shared.RuntimeValue, *errors.RuntimeError) {
	if lhs.Type != shared.Number || rhs.Type != shared.Number {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("This binary operation can only be performed on numbers. Attempted to perform division on `%s` and `%s`", shared.Stringify(lhs.Type), shared.Stringify(rhs.Type)),
		}
	}

	if rhs.Value.(float64) == 0 {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot perform division by zero. Attempted to divide `%g` by `0`", lhs.Value.(float64)),
		}
	}

	result := values.MK_NUMBER(0)

	lhsValue := lhs.Value.(float64)
	rhsValue := rhs.Value.(float64)

	lhsDecimalPlaces := countDecimalPlaces(lhsValue)
	rhsDecimalPlaces := countDecimalPlaces(rhsValue)

	scalingFactor := math.Pow10(max(lhsDecimalPlaces, rhsDecimalPlaces))

	lhsValue = lhsValue * float64(scalingFactor)
	rhsValue = rhsValue * float64(scalingFactor)

	r := float64(0)
	r = lhsValue / rhsValue

	result = values.MK_NUMBER(r)

	return &result, nil
}
