package debugger

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/environment"
)

type State string

const (
	RunningState  State = "running"
	PausedState   State = "paused"
	SteppingState State = "stepping"
)

type Debugger struct {
	BreakpointManager BreakpointManager
	Environment       *environment.Environment
	State             State
}

func NewDebugger(env *environment.Environment) *Debugger {
	return &Debugger{
		BreakpointManager: *NewBreakpointManager(),
		Environment:       env,
		State:             RunningState,
	}
}

// Define what can be debugged
var Debuggables = map[ast.NodeType]struct{}{
	ast.VarDeclarationNode:    {},
	ast.VarAssignmentExprNode: {},
	ast.IfStatementNode:       {},
	ast.WhileLoopNode:         {},
	ast.ReturnStmtNode:        {},
	ast.ContinueStmtNode:      {},
	ast.BreakStmtNode:         {},
	ast.TryCatchStmtNode:      {},
	ast.CallExprNode:          {},
	ast.FnDeclarationNode:     {},
	ast.ClassNode:             {},
	ast.ClassMethodNode:       {},
	ast.ClassPropertyNode:     {},
	ast.ProgramNode:           {},
}

// Internal API

func (d *Debugger) ShouldStop(filename string, line int) bool {
	return d.BreakpointManager.Has(filename, line)
}

func (d *Debugger) IsDebuggable(nodeType ast.NodeType) bool {
	_, isDebuggable := Debuggables[nodeType]
	return isDebuggable
}

// End user API

func (d *Debugger) Continue() error {
	d.State = RunningState
	return nil
}

func (d *Debugger) Step() error {
	d.State = SteppingState
	return nil
}

func (d *Debugger) StepOut() error {
	return nil
}

func (d *Debugger) StepInto() error {
	return nil
}

func (d *Debugger) Run() error {
	d.State = RunningState
	return nil
}

func (d *Debugger) Pause() error {
	d.State = PausedState
	return nil
}

func (d *Debugger) Stop() error {
	d.State = RunningState
	return nil
}
