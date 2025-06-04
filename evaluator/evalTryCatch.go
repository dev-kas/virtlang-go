package evaluator

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalTryCatch(node *ast.TryCatchStmt, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	scope := environment.NewEnvironment(env)

	// I realized that just removing the snapshot at catching the error
	// is not enough, we have to make sure the call stack has the same
	// length as things go in and out of this function, so that the snapshots
	// isnt messed up. snapshots are meant to be taken only upon uncaught errors
	// and not on caught errors
	var stackDepth int
	if dbgr != nil {
		stackDepth = len(dbgr.CallStack)
	}

	for _, stmt := range node.Try {

		_, err := Evaluate(stmt, scope, dbgr)
		if err != nil {
			// save the last snapshot
			var lastSnapshot debugger.Snapshot
			if dbgr != nil && len(dbgr.Snapshots) > 0 {
				lastSnapshot = dbgr.Snapshots[len(dbgr.Snapshots)-1]
			}

			// Because we are successfully catching this error,
			// we can safely delete the snapshot
			// Here we will remove everything other than the stackDepth
			// we saved earlier
			if dbgr != nil && len(dbgr.Snapshots) > 0 {
				dbgr.Snapshots = dbgr.Snapshots[:stackDepth]
			}
			scope = environment.NewEnvironment(env)

			// declare the catch variable
			catchVar_message := values.MK_STRING(err.Message)
			catchVar_stack_raw := []shared.RuntimeValue{}
			for _, frame := range lastSnapshot.Stack {
				name := values.MK_STRING(frame.Name)
				line := values.MK_NUMBER(float64(frame.Line))
				file := values.MK_STRING(frame.Filename)
				catchVar_stack_raw = append(catchVar_stack_raw, values.MK_OBJECT(map[string]*shared.RuntimeValue{
					"name": &name,
					"line": &line,
					"file": &file,
				}))
			}
			catchVar_stack := values.MK_ARRAY(catchVar_stack_raw)
			catchVar := values.MK_OBJECT(map[string]*shared.RuntimeValue{
				"message": &catchVar_message,
				"stack":   &catchVar_stack,
			})
			_, err = scope.DeclareVar(node.CatchVar, catchVar, false)
			if err != nil {
				return nil, err
			}

			var lastResult *shared.RuntimeValue = nil
			for _, stmt := range node.Catch {
				res, catchErr := Evaluate(stmt, scope, dbgr)
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
