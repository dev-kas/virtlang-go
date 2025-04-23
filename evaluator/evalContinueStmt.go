package evaluator

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

func evalContinueStmt(_ *ast.ContinueStmt, _ *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	retValue := errors.RuntimeError{
		Message: "`continue` statement used outside of a loop context.",
		InternalCommunicationProtocol: &errors.InternalCommunicationProtocol{
			Type: errors.ICP_Continue,
		},
	}

	val := values.MK_NIL()
	return &val, &retValue
}
