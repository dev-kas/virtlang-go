package evaluator

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/ast"
	"github.com/dev-kas/virtlang-go/environment"
	"github.com/dev-kas/virtlang-go/errors"
	"github.com/dev-kas/virtlang-go/shared"
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
