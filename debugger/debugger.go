package debugger

import (
	"sync"

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
	CurrentFile       string
	CurrentLine       int
	mu                sync.Mutex
	cond              *sync.Cond
}

func NewDebugger(env *environment.Environment) *Debugger {
	d := &Debugger{
		BreakpointManager: *NewBreakpointManager(),
		Environment:       env,
		State:             RunningState,
		mu:                sync.Mutex{},
	}
	d.cond = sync.NewCond(&d.mu)
	return d
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

func (d *Debugger) WaitIfPaused() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for d.State == PausedState {
		d.cond.Wait()
	}
	if d.State == SteppingState {
		d.State = PausedState
	}
}

// End user API

func (d *Debugger) Continue() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = RunningState
	d.cond.Broadcast()
	return nil
}

func (d *Debugger) Step() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = SteppingState
	d.cond.Broadcast()
	return nil
}

func (d *Debugger) StepOut() error {
	// TODO: Handle stepping out
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = SteppingState
	d.cond.Broadcast()
	return nil
}

func (d *Debugger) StepInto() error {
	// TODO: Handle stepping into
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = SteppingState
	d.cond.Broadcast()
	return nil
}

func (d *Debugger) Pause() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = PausedState
	return nil
}
