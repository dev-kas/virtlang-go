package evaluator

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

func evalClass(node *ast.Class, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	class := values.MK_CLASS(node.Name, node.Body, node.Constructor, env)

	return env.DeclareVar(node.Name, class, true)
}
