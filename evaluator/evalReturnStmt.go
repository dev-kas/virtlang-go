package evaluator

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/debugger"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

func evalReturnStmt(node *ast.ReturnStmt, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	evaluated, err := Evaluate(node.Value, env, dbgr)
	if err != nil {
		return nil, err
	}

	retValue := errors.RuntimeError{
		Message: "<RETURN STATEMENT>",
		InternalCommunicationProtocol: &errors.InternalCommunicationProtocol{
			Type:   errors.ICP_Return,
			RValue: evaluated,
		},
	}

	val := values.MK_NIL()
	return &val, &retValue
}
