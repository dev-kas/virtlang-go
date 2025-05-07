package evaluator

import (
	"github.com/dev-kas/virtlang-go/ast"
	"github.com/dev-kas/virtlang-go/environment"
	"github.com/dev-kas/virtlang-go/errors"
	"github.com/dev-kas/virtlang-go/shared"
	"github.com/dev-kas/virtlang-go/values"
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
