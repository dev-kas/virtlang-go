package evaluator

import (
	"github.com/dev-kas/virtlang-go/ast"
	"github.com/dev-kas/virtlang-go/environment"
	"github.com/dev-kas/virtlang-go/errors"
	"github.com/dev-kas/virtlang-go/shared"
	"github.com/dev-kas/virtlang-go/values"
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

		if evaluated != nil {
			result = *evaluated
		}
	}

	return &result, nil
}
