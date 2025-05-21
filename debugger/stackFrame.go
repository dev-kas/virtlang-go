package debugger

import "github.com/dev-kas/virtlang-go/v3/environment"

type StackFrame struct {
	Name     string
	Filename string
	Line     int
}

type CallStack []StackFrame

type Snapshot struct {
	Stack CallStack
	Env *environment.Environment
}

type Snapshots []Snapshot
