package evaluator

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
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
