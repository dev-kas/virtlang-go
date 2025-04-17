package evaluator

import (
	"VirtLang/ast"
	"VirtLang/environment"
	"VirtLang/errors"
	"VirtLang/shared"
)

func evalObjectExpr(o *ast.ObjectLiteral, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	obj := &shared.RuntimeValue{
		Type:  shared.Object,
		Value: map[string]*shared.RuntimeValue{},
	}

	for _, property := range o.Properties {
		var runtimeVal shared.RuntimeValue

		if property.Value == nil {
			val, err := env.LookupVar(property.Key)
			if err != nil {
				return nil, err
			}

			runtimeVal = *val
		} else {
			val, err := Evaluate(property.Value, env)
			if err != nil {
				return nil, err
			}

			runtimeVal = *val
		}

		obj.Value.(map[string]*shared.RuntimeValue)[property.Key] = &runtimeVal
	}

	return obj, nil
}
