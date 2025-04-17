package evaluator

import (
	"fmt"

	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
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
		nativeFn, err := fn.Value.(values.NativeFunctionValue)
		if !err {
			return nil, &errors.RuntimeError{
				Message: fmt.Sprintf("Unable to resolve native function type: %v.", fn.Type), // TODO: Write a stringify function to make the types human readable
			}
		}
		convertedArgs := make([]shared.RuntimeValue, len(args))
		for i, arg := range args {
			convertedArgs[i] = *arg
		}
		result := nativeFn.Call(convertedArgs, env)
		return &result, nil
	} else if fn.Type == shared.Function {
		fnVal := fn.Value.(values.FunctionValue)
		scope := environment.NewEnvironment(fnVal.DeclarationEnv)

		for i, param := range fnVal.Params {
			// TODO: check bounts and arity of fn
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
	} else {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot invoke a non-function (attempted to call a %v).", fn.Type), // TODO: Write a stringify function to make the types human readable
		}
	}
}
