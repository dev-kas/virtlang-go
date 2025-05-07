package evaluator

import (
	"github.com/dev-kas/virtlang-go/ast"
	"github.com/dev-kas/virtlang-go/environment"
	"github.com/dev-kas/virtlang-go/errors"
	"github.com/dev-kas/virtlang-go/shared"
	"github.com/dev-kas/virtlang-go/values"
)

func evalVarDecl(node *ast.VarDeclaration, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	var value shared.RuntimeValue

	if node.Value != nil {
		evaluated, err := Evaluate(node.Value, env)
		if err != nil {
			return nil, err
		}

		value = *evaluated
	} else {
		value = values.MK_NIL()
	}

	result, err := env.DeclareVar(node.Identifier, value, node.Constant)

	if err != nil {
		return nil, err
	}

	return result, nil
}
