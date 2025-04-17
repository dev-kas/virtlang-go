package evaluator

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
)

func evalFnDecl(node *ast.FnDeclaration, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	fn := values.FunctionValue{
		Name:           node.Name,
		Params:         node.Params,
		DeclarationEnv: env,
		Body:           node.Body,
		Type:           shared.Function,
		Value:          nil,
	}

	if node.Anonymous {
		return &shared.RuntimeValue{
			Type:  shared.Function,
			Value: fn,
		}, nil
	} else {
		return env.DeclareVar(node.Name, shared.RuntimeValue{Type: shared.Function, Value: fn}, true)
	}
}
