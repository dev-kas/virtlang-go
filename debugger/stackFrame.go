package debugger

type StackFrame struct {
	Name     string
	Filename string
	Line     int
}

type CallStack []StackFrame
