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

type StepType string

const (
	StepInto StepType = "step_into"
	StepOver StepType = "step_over"
	StepOut  StepType = "step_out"
)

type Debugger struct {
	BreakpointManager BreakpointManager
	Environment       *environment.Environment
	State             State
	CurrentFile       string
	CurrentLine       int
	CallStack         CallStack
	Snapshots         Snapshots
	stepType          StepType
	stepDepth         int
	mu                sync.Mutex
	cond              *sync.Cond
}

func NewDebugger(env *environment.Environment) *Debugger {
	d := &Debugger{
		BreakpointManager: *NewBreakpointManager(),
		Environment:       env,
		State:             RunningState,
		CallStack:         make(CallStack, 0),
		Snapshots:         make(Snapshots, 0),
		stepType:          StepInto,
		stepDepth:         0,
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
	ast.ProgramNode:{},
}

// Internal API

func (d *Debugger) ShouldStop(filename string, line int) bool {
	return d.BreakpointManager.Has(filename, line)
}

func (d *Debugger) IsDebuggable(nodeType ast.NodeType) bool {
	_, isDebuggable := Debuggables[nodeType]
	return isDebuggable
}

func (d *Debugger) WaitIfPaused(nodeType ast.NodeType) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for {
		switch d.State {
		case PausedState:
			d.cond.Wait()
		case SteppingState:
			if !d.IsDebuggable(nodeType) {
				// Not debuggable, continue execution
				return
			}
			currentDepth := len(d.CallStack)
			switch d.stepType {
			case StepInto:
				d.State = PausedState
				return
			case StepOver:
				if currentDepth <= d.stepDepth {
					d.State = PausedState
					return
				}
			case StepOut:
				if currentDepth < d.stepDepth {
					d.State = PausedState
					return
				}
			}
			d.cond.Wait()
		default:
			return
		}
	}
}

func (d *Debugger) PushFrame(frame StackFrame) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.CallStack = append(d.CallStack, frame)
}

func (d *Debugger) PopFrame() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.CallStack) == 0 {
		return
	}
	d.CallStack = d.CallStack[:len(d.CallStack)-1]
}

func (d *Debugger) TakeSnapshot() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Create deep copies of the call stack and environment
	stackCopy := DeepCopyCallStack(d.CallStack)
	envCopy := environment.DeepCopy(d.Environment)

	d.Snapshots = append(d.Snapshots, Snapshot{
		Stack: stackCopy,
		Env:   envCopy,
	})
}

// DeepCopyCallStack creates a deep copy of the call stack
func DeepCopyCallStack(stack CallStack) CallStack {
	if stack == nil {
		return nil
	}

	newStack := make(CallStack, len(stack))
	copy(newStack, stack)
	return newStack
}

// End user API

func (d *Debugger) Continue() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = RunningState
	d.cond.Broadcast()
	return nil
}

func (d *Debugger) StepInto() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = SteppingState
	d.stepType = StepInto
	d.stepDepth = len(d.CallStack)
	d.cond.Broadcast()
	return nil
}

func (d *Debugger) StepOver() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = SteppingState
	d.stepType = StepOver
	d.stepDepth = len(d.CallStack)
	d.cond.Broadcast()
	return nil
}

func (d *Debugger) StepOut() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = SteppingState
	d.stepType = StepOut
	d.stepDepth = len(d.CallStack)
	d.cond.Broadcast()
	return nil
}

func (d *Debugger) Pause() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = PausedState
	return nil
}
