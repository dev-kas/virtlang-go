package values

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
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

func MK_OBJECT(value map[string]*shared.RuntimeValue) shared.RuntimeValue {
	properties := map[string]*shared.RuntimeValue{}

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

type ClassValue struct {
	Type           shared.ValueType
	Value          any
	Name           string
	Body           []ast.Stmt
	DeclarationEnv *environment.Environment
	Constructor    *ast.ClassMethod
}

func MK_CLASS(name string, body []ast.Stmt, constructor *ast.ClassMethod, declarationEnv *environment.Environment) shared.RuntimeValue {
	return shared.RuntimeValue{
		Type: shared.Class,
		Value: ClassValue{
			Type:           shared.Class,
			Value:          nil,
			Name:           name,
			Body:           body,
			DeclarationEnv: declarationEnv,
			Constructor:    constructor,
		},
	}
}

type ClassInstanceValue struct {
	Type    shared.ValueType
	Class   ClassValue
	Publics map[string]bool
	Data    *environment.Environment
}

func MK_CLASS_INSTANCE(class *ClassValue, publics map[string]bool, data *environment.Environment) shared.RuntimeValue {
	return shared.RuntimeValue{
		Type: shared.ClassInstance,
		Value: ClassInstanceValue{
			Type:    shared.ClassInstance,
			Class:   *class,
			Publics: publics,
			Data:    data,
		},
	}
}
