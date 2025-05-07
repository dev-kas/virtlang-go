package evaluator

import (
	"github.com/dev-kas/virtlang-go/ast"
	"github.com/dev-kas/virtlang-go/environment"
	"github.com/dev-kas/virtlang-go/errors"
	"github.com/dev-kas/virtlang-go/shared"
	"github.com/dev-kas/virtlang-go/values"
)

func evalReturnStmt(node *ast.ReturnStmt, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	evaluated, err := Evaluate(node.Value, env)
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
