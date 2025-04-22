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
	default:
		return "unknown"
	}
}
