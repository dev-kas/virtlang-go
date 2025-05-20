package evaluator

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

func evalCallExpr(node *ast.CallExpr, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	args := make([]*shared.RuntimeValue, len(node.Args))
	for i, arg := range node.Args {
		evaluatedArg, err := Evaluate(arg, env)
		if err != nil {
			return nil, err
		}
		args[i] = evaluatedArg
	}

	fn, err := Evaluate(node.Callee, env)
	if err != nil {
		return nil, err
	}

	if fn.Type == shared.NativeFN {
		nativeFn, err := fn.Value.(values.NativeFunction)
		if !err {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Unable to resolve native function type: %s.", shared.Stringify(fn.Type)),
			}
		}
		convertedArgs := make([]shared.RuntimeValue, len(args))
		for i, arg := range args {
			convertedArgs[i] = *arg
		}
		result, call_err := nativeFn(convertedArgs, env)
		if call_err != nil {
			return nil, call_err
		}
		return result, nil
	} else if fn.Type == shared.Function {
		fnVal := fn.Value.(values.FunctionValue)
		scope := environment.NewEnvironment(fnVal.DeclarationEnv)

		for i, param := range fnVal.Params {
			// TODO: check bounds and arity of fn
			scope.DeclareVar(param, *args[i], true)
		}

		result := values.MK_NIL()

		for _, stmt := range fnVal.Body {
			var err *errors.RuntimeError
			var res *shared.RuntimeValue
			res, err = Evaluate(stmt, &scope)
			if err != nil {
				if err.InternalCommunicationProtocol != nil && err.InternalCommunicationProtocol.Type == errors.ICP_Return {
					return err.InternalCommunicationProtocol.RValue, nil
				}
				return nil, err
			}
			result = *res
		}

		return &result, nil
	} else if fn.Type == shared.Class {
		classVal := fn.Value.(values.ClassValue)
		// forkedEnv := environment.NewEnvironment(classVal.DeclarationEnv)
		// instance := &values.ClassInstanceValue{
		// 	Class: &classVal,
		// 	Env:   &forkedEnv,
		// }

		if classVal.Constructor == nil {
			return nil, &errors.RuntimeError{
				Message: "Class has no constructor.",
			}
		}

		classScope := environment.NewEnvironment(classVal.DeclarationEnv)
		publics := map[string]bool{}
		for _, stmt := range classVal.Body {
			if stmt.GetType() == ast.ClassMethodNode {
				method := stmt.(*ast.ClassMethod)
				_, err := evalClassMethod(method, &classScope)
				if err != nil {
					return nil, err
				}
				if method.IsPublic {
					publics[method.Name] = true
				}
			} else if stmt.GetType() == ast.ClassPropertyNode {
				property := stmt.(*ast.ClassProperty)
				_, err := evalClassProperty(property, &classScope)
				if err != nil {
					return nil, err
				}
				if property.IsPublic {
					publics[property.Name] = true
				}
			}
		}

		constructor := classVal.Constructor
		constructorScope := environment.NewEnvironment(&classScope)
		for i, param := range constructor.Params {
			// TODO: check bounds and arity of fn
			constructorScope.DeclareVar(param, *args[i], true)
		}

		for _, stmt := range constructor.Body {
			_, err := Evaluate(stmt, &constructorScope)
			if err != nil {
				if err.InternalCommunicationProtocol != nil && err.InternalCommunicationProtocol.Type == errors.ICP_Return {
					return nil, &errors.RuntimeError{
						Message: "Constructor cannot return a value.",
					}
				}
				return nil, err
			}
		}

		retVal := values.MK_CLASS_INSTANCE(&classVal, publics, &classScope)
		return &retVal, nil
	} else {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot invoke a non-function (attempted to call a %s).", shared.Stringify(fn.Type)),
		}
	}
}
