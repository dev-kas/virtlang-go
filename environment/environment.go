package environment

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
)

type Environment struct {
	Parent    *Environment
	Variables map[string]shared.RuntimeValue
	Constants map[string]struct{}
	Global    bool
}

func NewEnvironment(fork *Environment) Environment {
	global := false
	if fork == nil {
		global = true
	}
	env := Environment{
		Parent:    fork,
		Variables: map[string]shared.RuntimeValue{},
		Constants: map[string]struct{}{},
		Global:    global,
	}

	return env
}

func (e *Environment) DeclareVar(name string, value shared.RuntimeValue, constant bool) (*shared.RuntimeValue, *errors.RuntimeError) {
	if name == "" { // Trimming space is intentionally omitted here
		return &value, nil
	}

	if _, exists := e.Variables[name]; exists {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot redeclare variable `%s`", name),
		}
	}

	e.Variables[name] = value
	if constant {
		e.Constants[name] = struct{}{}
	}

	return &value, nil
}

func (e *Environment) Resolve(varname string) (*Environment, *errors.RuntimeError) {
	if _, exists := e.Variables[varname]; exists {
		return e, nil
	}

	if e.Parent == nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot resolve variable `%s`", varname),
		}
	}

	return e.Parent.Resolve(varname)
}

func (e *Environment) LookupVar(name string) (*shared.RuntimeValue, *errors.RuntimeError) {
	env, err := e.Resolve(name)
	if err != nil {
		return nil, err
	}
	value := env.Variables[name]

	if value.Type == shared.Object || value.Type == shared.Array {
	}

	return &value, err
}

func (e *Environment) AssignVar(name string, value shared.RuntimeValue) (*shared.RuntimeValue, *errors.RuntimeError) {
	if name == "" { // Trimming space is intentionally omitted here
		return &value, nil
	}

	env, err := e.Resolve(name)
	if err != nil {
		return nil, err
	}
	if _, exists := env.Constants[name]; exists {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot reassign to constant variable `%s`", name),
		}
	}

	if value.Type == shared.Array || value.Type == shared.Object {
		currentValue, exists := env.Variables[name]
		if exists && (currentValue.Type == shared.Array || currentValue.Type == shared.Object) {
			currentValue.Value = value.Value
			env.Variables[name] = currentValue
		} else {
			env.Variables[name] = value
		}
	} else {
		env.Variables[name] = value
	}

	return &value, nil
}
