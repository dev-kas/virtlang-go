package evaluator

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/debugger"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/shared"
	"github.com/dev-kas/virtlang-go/v3/values"
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

		_, err := Evaluate(stmt, &scope, dbgr)
		if err != nil {
			// Because we are successfully catching this error,
			// we can safely delete the snapshot
			// Here we will remove everything other than the stackDepth
			// we saved earlier
			if dbgr != nil && len(dbgr.Snapshots) > 0 {
				dbgr.Snapshots = dbgr.Snapshots[:stackDepth]
			}
			scope = environment.NewEnvironment(env)

			scope.DeclareVar(node.CatchVar, shared.RuntimeValue{Type: shared.String, Value: "\"Runtime Error: " + err.Message + "\""}, false)

			var lastResult *shared.RuntimeValue = nil
			for _, stmt := range node.Catch {
				res, catchErr := Evaluate(stmt, &scope, dbgr)
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
