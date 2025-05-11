package shared

type ValueType int

const (
	Nil ValueType = iota
	Number
	Boolean
	Object
	Array
	NativeFN
	Function
	String
	Class
	ClassInstance
)

type RuntimeValue struct {
	Type  ValueType
	Value any
}

func Stringify(v ValueType) string {
	switch v {
	case Nil:
		return "nil"
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	case Object:
		return "object"
	case Array:
		return "array"
	case NativeFN:
		return "native-function"
	case Function:
		return "function"
	case String:
		return "string"
	case Class:
		return "class"
	case ClassInstance:
		return "class-instance"
	default:
		return "unknown"
	}
}
