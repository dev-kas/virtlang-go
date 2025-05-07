package values

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
)

// Factory Functions

func MK_NIL() shared.RuntimeValue {
	return shared.RuntimeValue{
		Type:  shared.Nil,
		Value: nil,
	}
}

func MK_BOOL(value bool) shared.RuntimeValue {
	return shared.RuntimeValue{
		Type:  shared.Boolean,
		Value: value,
	}
}

func MK_NUMBER(value float64) shared.RuntimeValue {
	return shared.RuntimeValue{
		Type:  shared.Number,
		Value: value,
	}
}

func MK_STRING(value string) shared.RuntimeValue {
	return shared.RuntimeValue{
		Type:  shared.String,
		Value: value,
	}
}

func MK_OBJECT(value map[string]shared.RuntimeValue) shared.RuntimeValue {
	properties := map[string]shared.RuntimeValue{}

	for key, val := range value {
		properties[key] = val
	}

	return shared.RuntimeValue{
		Type:  shared.Object,
		Value: properties,
	}
}

type NativeFunction func(args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError)

func MK_NATIVE_FN(fn NativeFunction) shared.RuntimeValue {
	return shared.RuntimeValue{
		Type:  shared.NativeFN,
		Value: fn,
	}
}

type FunctionValue struct {
	Type           shared.ValueType
	Value          any
	Name           string
	Params         []string
	DeclarationEnv *environment.Environment
	Body           []ast.Stmt
}

func MK_ARRAY(value []shared.RuntimeValue) shared.RuntimeValue {
	return shared.RuntimeValue{
		Type:  shared.Array,
		Value: value,
	}
}
