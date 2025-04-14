package values

import (
	"VirtLang/environment"
	"VirtLang/shared"
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

type FunctionCall func(args []shared.RuntimeValue, env *environment.Environment) shared.RuntimeValue
type NativeFunctionValue struct {
	Call  FunctionCall
	Type  shared.ValueType
	Value interface{}
}

func MK_NATIVE_FN(fn FunctionCall) NativeFunctionValue {
	return NativeFunctionValue{
		Type:  shared.NativeFN,
		Value: nil,
		Call:  fn,
	}
}
