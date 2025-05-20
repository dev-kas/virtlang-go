package evaluator

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
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
