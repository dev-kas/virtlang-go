package debugger

import "fmt"

type BreakpointManager struct {
	Breakpoints map[string]bool
}

func NewBreakpointManager() *BreakpointManager {
	return &BreakpointManager{
		Breakpoints: make(map[string]bool),
	}
}

func (bm *BreakpointManager) Set(file string, line int) {
	name := fmt.Sprintf("%s:%d", file, line)
	bm.Breakpoints[name] = true
}

func (bm *BreakpointManager) Remove(file string, line int) {
	name := fmt.Sprintf("%s:%d", file, line)
	delete(bm.Breakpoints, name)
}

func (bm *BreakpointManager) Clear() {
	bm.Breakpoints = make(map[string]bool)
}

func (bm *BreakpointManager) Has(file string, line int) bool {
	name := fmt.Sprintf("%s:%d", file, line)
	return bm.Breakpoints[name]
}
