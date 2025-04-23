package evaluator

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

func evalBreakStmt(_ *ast.BreakStmt, _ *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	retValue := errors.RuntimeError{
		Message: "`break` statement used outside of a loop context.",
		InternalCommunicationProtocol: &errors.InternalCommunicationProtocol{
			Type: errors.ICP_Break,
		},
	}

	val := values.MK_NIL()
	return &val, &retValue
}
