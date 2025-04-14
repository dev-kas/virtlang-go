package shared

type ValueType int

const (
	Nil ValueType = iota
	Number
	Boolean
	Object
	NativeFN
	Function
	String
)

type RuntimeValue struct {
	Type  ValueType
	Value any
}
