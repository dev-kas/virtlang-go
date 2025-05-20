package evaluator

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/helpers"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

func evalIfStmt(statement *ast.IfStatement, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	cond, err := Evaluate(statement.Condition, env)
	if err != nil {
		return nil, err
	}

	// Main `if` branch
	if helpers.IsTruthy(cond) {
		for _, stmt := range statement.Body {
			if _, err := Evaluate(stmt, env); err != nil {
				return nil, err
			}
		}
		nilVal := values.MK_NIL()
		return &nilVal, nil
	}

	// Else-if branches
	if len(statement.ElseIf) > 0 {
		// We only evaluate the first `else if` branch in the array
		// this is because I observed the parser that it does not create multiple
		// elements in array, but just one array that's deeply nested
		// and i thought it;s already a good option,
		// so we just evaluate the first element
		return evalIfStmt(statement.ElseIf[0], env)
	}

	// Else branch
	for _, stmt := range statement.Else {
		if _, err := Evaluate(stmt, env); err != nil {
			return nil, err
		}
	}

	nilVal := values.MK_NIL()
	return &nilVal, nil
}
