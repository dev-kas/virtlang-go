package evaluator

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
)

func evalIfStmt(statement *ast.IfStatement, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	ifStmt := ast.IfStatement{
		Body:      statement.Body,
		Condition: statement.Condition,
	}

	conditionSatisfied, err := Evaluate(ifStmt.Condition, env)
	if err != nil {
		return nil, err
	}

	if conditionSatisfied.Type == shared.Boolean && conditionSatisfied.Value.(bool) {
		for _, stmt := range ifStmt.Body {
			_, err := Evaluate(stmt, env)
			if err != nil {
				return nil, err
			}
		}
	}

	retValue := values.MK_NIL()
	return &retValue, nil
}
