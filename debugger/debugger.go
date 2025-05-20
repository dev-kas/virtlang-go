package debugger

import (
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

// Internal API

func (d *Debugger) ShouldStop(line string, col int) bool {
	return d.BreakpointManager.Has(line, col)
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
