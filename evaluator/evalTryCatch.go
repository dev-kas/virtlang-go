package evaluator

import (
	"VirtLang/ast"
	"VirtLang/environment"
	"VirtLang/errors"
	"VirtLang/shared"
	"VirtLang/values"
)

func evalTryCatch(node *ast.TryCatchStmt, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	scope := environment.NewEnvironment(env)

	for _, stmt := range node.Try {
		_, err := Evaluate(stmt, &scope)
		if err != nil {
			scope = environment.NewEnvironment(env)

			scope.DeclareVar(node.CatchVar, shared.RuntimeValue{Type: shared.String, Value: err.Error()}, false)

			for _, stmt := range node.Catch {
				_, err := Evaluate(stmt, &scope)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	result := values.MK_NIL()
	return &result, nil
}