package values

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/shared"
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

func MK_NUMBER(value int) shared.RuntimeValue { // TODO: convert to float64 later
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

type NativeFunctionCall func(args []shared.RuntimeValue, env *environment.Environment) shared.RuntimeValue
type NativeFunctionValue struct {
	Call  NativeFunctionCall
	Type  shared.ValueType
	Value interface{}
}

func MK_NATIVE_FN(fn NativeFunctionCall) NativeFunctionValue {
	return NativeFunctionValue{
		Type:  shared.NativeFN,
		Value: nil,
		Call:  fn,
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
