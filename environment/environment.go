package environment

import (
	"VirtLang/errors"
	"VirtLang/shared"
	"fmt"
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
	return Environment{
		Parent:    fork,
		Variables: map[string]shared.RuntimeValue{},
		Constants: map[string]struct{}{},
		Global:    global,
	}
}

func (e *Environment) DeclareVar(name string, value shared.RuntimeValue, constant bool) (*shared.RuntimeValue, *errors.RuntimeError) {
	if name == "" { // Trimming space is intentionally omitted here
		return &value, nil // Currently skip, as a replication of the original code.
	} // TODO: Fix that above code later to see why the interpreter sometimes give an empty name.

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
	return &value, err
}

func (e *Environment) AssignVar(name string, value shared.RuntimeValue) (*shared.RuntimeValue, *errors.RuntimeError) {
	if name == "" { // Trimming space is intentionally omitted here
		return &value, nil // Currently skip, as a replication of the original code.
	} // TODO: Fix that above code later to see why the interpreter sometimes give an empty name.

	if _, exists := e.Variables[name]; !exists {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot reassign to constant variable `%s`", name),
		}
	}

	e.Variables[name] = value
	return &value, nil
}
