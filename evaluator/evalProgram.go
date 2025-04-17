package evaluator

import (
	"VirtLang/ast"
	"VirtLang/environment"
	"VirtLang/errors"
	"VirtLang/shared"
	"VirtLang/values"
)

func evalProgram(astNode *ast.Program, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	result := values.MK_NIL()

	for _, stmt := range astNode.Stmts {
		evaluated, err := Evaluate(stmt, env)
		if err != nil {
			if err.InternalCommunicationProtocol.Type == errors.ICP_Return {
				return err.InternalCommunicationProtocol.RValue, nil
			}
			return nil, err
		}

		result = *evaluated
	}

	return &result, nil
}
