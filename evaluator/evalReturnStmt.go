package evaluator

import (
	"VirtLang/ast"
	"VirtLang/environment"
	"VirtLang/errors"
	"VirtLang/shared"
	"VirtLang/values"
)

func evalReturnStmt(node *ast.ReturnStmt, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	evaluated, err := Evaluate(node.Value, env)
	if err != nil {
		return nil, err
	}

	retValue := errors.RuntimeError{
		Message: "<RETURN STATEMENT>",
		InternalCommunicationProtocol: &errors.InternalCommunicationProtocol{
			Type: errors.ICP_Return,
			RValue: evaluated,
		},
	}

	val := values.MK_NIL()
	return &val, &retValue
}
