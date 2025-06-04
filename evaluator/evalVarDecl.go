package evaluator

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalVarDecl(node *ast.VarDeclaration, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	if node.Value != nil {
		evaluated, err := Evaluate(node.Value, env, dbgr)
		if err != nil {
			return nil, err
		}
		return env.DeclareVar(node.Identifier, *evaluated, node.Constant)
	}

	value := values.MK_NIL()
	return env.DeclareVar(node.Identifier, value, node.Constant)
}
