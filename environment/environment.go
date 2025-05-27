package environment

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
)

type Environment struct {
	Parent    *Environment
	Variables map[string]*shared.RuntimeValue
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
		Variables: make(map[string]*shared.RuntimeValue),
		Constants: make(map[string]struct{}),
		Global:    global,
	}

	return env
}

func (e *Environment) DeclareVar(name string, value shared.RuntimeValue, constant bool) (*shared.RuntimeValue, *errors.RuntimeError) {
	if name == "" { // Trimming space is intentionally omitted here
		val := value
		return &val, nil
	}

	if _, exists := e.Variables[name]; exists {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot redeclare variable `%s`", name),
		}
	}

	val := value
	e.Variables[name] = &val
	if constant {
		e.Constants[name] = struct{}{}
	}

	return e.Variables[name], nil
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
	value, exists := env.Variables[name]
	if !exists {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Variable '%s' not found", name),
		}
	}
	return value, nil
}

func (e *Environment) AssignVar(name string, value shared.RuntimeValue) (*shared.RuntimeValue, *errors.RuntimeError) {
	if name == "" { // Trimming space is intentionally omitted here
		val := value
		return &val, nil
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

	env.Variables[name] = &value
	return env.Variables[name], nil
}
