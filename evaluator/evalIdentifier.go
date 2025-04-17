package evaluator

import (
	"VirtLang/ast"
	"VirtLang/environment"
	"VirtLang/errors"
	"VirtLang/shared"
	"fmt"
)

func evalIdentifier(identifier *ast.Identifier, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	value, err := env.LookupVar(identifier.Symbol)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Undefined variable: %s", identifier.Symbol),
		}
	}
	return value, nil
}
