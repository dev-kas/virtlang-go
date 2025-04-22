package evaluator

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

func evalProgram(astNode *ast.Program, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	result := values.MK_NIL()

	for _, stmt := range astNode.Stmts {
		evaluated, err := Evaluate(stmt, env)
		if err != nil {
			if err.InternalCommunicationProtocol != nil && err.InternalCommunicationProtocol.Type == errors.ICP_Return {
				return err.InternalCommunicationProtocol.RValue, nil
			}
			return nil, err
		}

		result = *evaluated
	}

	return &result, nil
}
