package evaluator

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
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
