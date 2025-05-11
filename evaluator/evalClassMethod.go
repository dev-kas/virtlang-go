package evaluator

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
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

	_, err := env.DeclareVar(node.Name, shared.RuntimeValue{
		Type:  shared.Function,
		Value: *method,
	}, false)

	if err != nil {
		return nil, err
	}

	result := shared.RuntimeValue{
		Type:  shared.Function,
		Value: *method,
	}

	return &result, nil
}
