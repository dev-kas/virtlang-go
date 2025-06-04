package evaluator

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalProgram(astNode *ast.Program, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	// Push the <main> frame
	if dbgr != nil && dbgr.IsDebuggable(ast.ProgramNode) {
		dbgr.PushFrame(debugger.StackFrame{
			Name:     "<main>",
			Filename: astNode.GetSourceMetadata().Filename,
			Line:     astNode.GetSourceMetadata().StartLine,
		})
		defer dbgr.PopFrame()
	}

	result := values.MK_NIL()

	for _, stmt := range astNode.Stmts {
		evaluated, err := Evaluate(stmt, env, dbgr)
		if err != nil {
			if err.InternalCommunicationProtocol != nil && err.InternalCommunicationProtocol.Type == errors.ICP_Return {
				return err.InternalCommunicationProtocol.RValue, nil
			}
			// Take snapshot
			if dbgr != nil {
				dbgr.TakeSnapshot()
			}
			return nil, err
		}

		if evaluated != nil {
			result = *evaluated
		}
	}

	return &result, nil
}
