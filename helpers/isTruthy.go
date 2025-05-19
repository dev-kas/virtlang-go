package helpers

import (
	"github.com/dev-kas/virtlang-go/v2/shared"
)

// IsTruthy determines whether a VirtLang RuntimeValue should be considered
// "truthy" in boolean contexts (like if statements and while loops).
//
// Truthiness rules:
// - Boolean: true is truthy, false is falsy
// - Number: non-zero is truthy, zero is falsy
// - String: non-empty is truthy, empty is falsy
// - Nil: always falsy
// - Object: always truthy (even empty objects)
// - Array: always truthy (even empty arrays)
// - Function: always truthy
// - NativeFN: always truthy
// - ClassInstance: always truthy
// - Class: always truthy
// - Unknown: always truthy
func IsTruthy(value *shared.RuntimeValue) bool {
	if value == nil {
		return false
	}

	switch value.Type {
	case shared.Boolean:
		// Boolean values are truthy if they're true
		return value.Value.(bool)

	case shared.Number:
		// Numbers are truthy if they're non-zero
		num := value.Value.(float64)
		return num != 0

	case shared.String:
		// Strings are truthy if they're non-empty
		str := value.Value.(string)
		// We are using 2 here because the string, specifically VirtLang's string RuntimeValue,
		// is wrapped in double quotes which adds 2 characters to the string.
		return len(str) > 2

	case shared.Nil:
		// nil is always falsy
		return false

	case shared.Object, shared.Array, shared.Function, shared.NativeFN, shared.ClassInstance, shared.Class:
		// Objects, arrays, and functions are always truthy
		return true

	default:
		// For any unknown types, default to truthy
		// I might have to reconsider my decision here
		// but it's fine as long as we maintain this file
		// if we ever add new types, which, is kinda rare because
		// VirtLang already has a good set of types
		return true
	}
}
