package evaluator

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

func evalTryCatch(node *ast.TryCatchStmt, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	scope := environment.NewEnvironment(env)

	for _, stmt := range node.Try {
		_, err := Evaluate(stmt, &scope)
		if err != nil {
			scope = environment.NewEnvironment(env)

			scope.DeclareVar(node.CatchVar, shared.RuntimeValue{Type: shared.String, Value: "Runtime Error: " + err.Message}, false)

			var lastResult *shared.RuntimeValue = nil
			for _, stmt := range node.Catch {
				res, catchErr := Evaluate(stmt, &scope)
				if catchErr != nil {
					return nil, catchErr
				}
				lastResult = res
			}
			if lastResult != nil {
				return lastResult, nil
			}
			result := values.MK_NIL()
			return &result, nil
		}
	}

	result := values.MK_NIL()
	return &result, nil
}
