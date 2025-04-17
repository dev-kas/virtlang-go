package evaluator

import (
	"VirtLang/ast"
	"VirtLang/environment"
	"VirtLang/errors"
	"VirtLang/shared"
	"VirtLang/values"
)

func evalWhileLoop(astNode *ast.WhileLoop, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	for {
		cond, err := Evaluate(astNode.Condition, env)
		if err != nil {
			return nil, err
		}

		if cond.Type == shared.Boolean && cond.Value.(bool) {
			scope := environment.NewEnvironment(env)
			for _, stmt := range astNode.Body {
				_, err := Evaluate(stmt, &scope)
				if err != nil {
					return nil, err
				}
			}
		} else {
			break
		}
	}
	
	result := values.MK_NIL()
	return &result, nil
}
