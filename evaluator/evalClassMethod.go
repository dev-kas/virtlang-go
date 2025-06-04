package evaluator

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalClassMethod(node *ast.ClassMethod, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	method := &values.FunctionValue{
		Name:           node.Name,
		Params:         node.Params,
		DeclarationEnv: env,
		Body:           node.Body,
		Type:           shared.Function,
		Value:          nil,
	}

	rtValue := &shared.RuntimeValue{
		Type:  shared.Function,
		Value: method,
	}

	_, err := env.DeclareVar(node.Name, *rtValue, false)
	if err != nil {
		return nil, err
	}

	return rtValue, nil
}
