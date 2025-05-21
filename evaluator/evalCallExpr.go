package evaluator

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/debugger"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
)

func evalCallExpr(node *ast.CallExpr, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	args := make([]*shared.RuntimeValue, len(node.Args))
	for i, arg := range node.Args {
		evaluatedArg, err := Evaluate(arg, env, dbgr)
		if err != nil {
			return nil, err
		}
		args[i] = evaluatedArg
	}

	fn, err := Evaluate(node.Callee, env, dbgr)
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

		// Push frame to stack
		if dbgr != nil {
			dbgr.PushFrame(debugger.StackFrame{
				Name: func() string {
					if fnVal.Name == "" {
						return "<anonymous>"
					}
					return fnVal.Name
				}(),
				Filename: node.GetSourceMetadata().Filename,
				Line:     node.GetSourceMetadata().StartLine,
			})
		}

		for i, param := range fnVal.Params {
			// TODO: check bounds and arity of fn
			scope.DeclareVar(param, *args[i], true)
		}

		result := values.MK_NIL()

		for _, stmt := range fnVal.Body {
			var err *errors.RuntimeError
			var res *shared.RuntimeValue
			res, err = Evaluate(stmt, &scope, dbgr)
			if err != nil {
				// pop frame from stack
				if dbgr != nil {
					dbgr.PopFrame()
				}
				if err.InternalCommunicationProtocol != nil && err.InternalCommunicationProtocol.Type == errors.ICP_Return {
					return err.InternalCommunicationProtocol.RValue, nil
				}
				// Take a snapshot now because it is a real error,
				// not a control flow event
				if dbgr != nil {
					dbgr.TakeSnapshot()
				}
				return nil, err
			}
			result = *res
		}

		// Pop frame from stack
		if dbgr != nil {
			dbgr.PopFrame()
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
				_, err := evalClassProperty(property, &classScope, dbgr)
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

		// Push frame to stack
		if dbgr != nil {
			dbgr.PushFrame(debugger.StackFrame{
				Name:     "constructor",
				Filename: node.GetSourceMetadata().Filename,
				Line:     node.GetSourceMetadata().StartLine,
			})
		}

		for _, stmt := range constructor.Body {
			_, err := Evaluate(stmt, &constructorScope, dbgr)
			if err != nil {
				// Take snapshot and pop frame from stack
				if dbgr != nil {
					dbgr.TakeSnapshot()
					dbgr.PopFrame()
				}
				if err.InternalCommunicationProtocol != nil && err.InternalCommunicationProtocol.Type == errors.ICP_Return {
					return nil, &errors.RuntimeError{
						Message: "Constructor cannot return a value.",
					}
				}
				return nil, err
			}
		}

		// Pop frame from stack
		if dbgr != nil {
			dbgr.PopFrame()
		}

		retVal := values.MK_CLASS_INSTANCE(&classVal, publics, &classScope)
		return &retVal, nil
	} else {
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Cannot invoke a non-function (attempted to call a %s).", shared.Stringify(fn.Type)),
		}
	}
}
