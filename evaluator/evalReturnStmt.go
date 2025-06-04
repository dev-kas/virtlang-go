package evaluator

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
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
