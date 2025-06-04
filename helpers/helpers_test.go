package helpers

import (
	"fmt"
	"testing"

	"github.com/dev-kas/virtlang-go/v3/shared"
)

func TestIsTruthy(t *testing.T) {
	testCases := []struct {
		name     string
		value    *shared.RuntimeValue
		expected bool
	}{
		// 1. Nil *shared.RuntimeValue (top-level nil pointer)
		{
			name:     "NilPointerInput",
			value:    nil,
			expected: false,
		},

		// 2. Boolean Type
		{
			name:     "BooleanTrue",
			value:    &shared.RuntimeValue{Type: shared.Boolean, Value: true},
			expected: true,
		},
		{
			name:     "BooleanFalse",
			value:    &shared.RuntimeValue{Type: shared.Boolean, Value: false},
			expected: false,
		},

		// 3. Number Type
		{
			name:     "NumberPositive",
			value:    &shared.RuntimeValue{Type: shared.Number, Value: float64(10.5)},
			expected: true,
		},
		{
			name:     "NumberNegative",
			value:    &shared.RuntimeValue{Type: shared.Number, Value: float64(-5)},
			expected: true,
		},
		{
			name:     "NumberZero",
			value:    &shared.RuntimeValue{Type: shared.Number, Value: float64(0.0)},
			expected: false,
		},
		{
			name:     "NumberExactlyZeroInteger",
			value:    &shared.RuntimeValue{Type: shared.Number, Value: float64(0)},
			expected: false,
		},

		// 4. String Type
		{
			name:     "StringNonEmptySimple",
			value:    &shared.RuntimeValue{Type: shared.String, Value: "a"},
			expected: true,
		},
		{
			name:     "StringNonEmptyLong",
			value:    &shared.RuntimeValue{Type: shared.String, Value: "hello world"},
			expected: true,
		},
		{
			name:     "StringEmptyVirtLang",
			value:    &shared.RuntimeValue{Type: shared.String, Value: ""},
			expected: false,
		},
		{
			name:     "StringWithSpaces",
			value:    &shared.RuntimeValue{Type: shared.String, Value: "   "},
			expected: true,
		},

		// 5. Nil Type (VirtLang's nil type, not a nil pointer)
		{
			name:     "NilType",
			value:    &shared.RuntimeValue{Type: shared.Nil, Value: nil},
			expected: false,
		},
		{
			name:     "NilTypeWithGarbageValue",
			value:    &shared.RuntimeValue{Type: shared.Nil, Value: "some garbage"},
			expected: false,
		},

		// 6. Object Type (always truthy)
		{
			name:     "ObjectEmpty",
			value:    &shared.RuntimeValue{Type: shared.Object, Value: make(map[string]*shared.RuntimeValue)},
			expected: true,
		},
		{
			name:     "ObjectNonEmpty",
			value:    &shared.RuntimeValue{Type: shared.Object, Value: map[string]*shared.RuntimeValue{"key": {Type: shared.Number, Value: float64(1)}}},
			expected: true,
		},
		{
			name:     "ObjectNilValueField",
			value:    &shared.RuntimeValue{Type: shared.Object, Value: nil},
			expected: true,
		},

		// 7. Array Type (always truthy)
		{
			name:     "ArrayEmpty",
			value:    &shared.RuntimeValue{Type: shared.Array, Value: make([]*shared.RuntimeValue, 0)},
			expected: true,
		},
		{
			name:     "ArrayNonEmpty",
			value:    &shared.RuntimeValue{Type: shared.Array, Value: []*shared.RuntimeValue{{Type: shared.Number, Value: float64(1)}}},
			expected: true,
		},
		{
			name:     "ArrayNilValueField",
			value:    &shared.RuntimeValue{Type: shared.Array, Value: nil},
			expected: true,
		},

		// 8. Function Type (always truthy)
		{
			name:     "FunctionType",
			value:    &shared.RuntimeValue{Type: shared.Function, Value: "function_representation"},
			expected: true,
		},
		{
			name:     "FunctionTypeNilValue",
			value:    &shared.RuntimeValue{Type: shared.Function, Value: nil},
			expected: true,
		},

		// 9. NativeFN Type (always truthy)
		{
			name:     "NativeFNType",
			value:    &shared.RuntimeValue{Type: shared.NativeFN, Value: "native_function_representation"},
			expected: true,
		},
		{
			name:     "NativeFNTypeNilValue",
			value:    &shared.RuntimeValue{Type: shared.NativeFN, Value: nil},
			expected: true,
		},

		// 10. ClassInstance Type (always truthy)
		{
			name:     "ClassInstanceType",
			value:    &shared.RuntimeValue{Type: shared.ClassInstance, Value: "instance_representation"},
			expected: true,
		},
		{
			name:     "ClassInstanceTypeNilValue",
			value:    &shared.RuntimeValue{Type: shared.ClassInstance, Value: nil},
			expected: true,
		},

		// 11. Class Type (always truthy)
		{
			name:     "ClassType",
			value:    &shared.RuntimeValue{Type: shared.Class, Value: "class_representation"},
			expected: true,
		},
		{
			name:     "ClassTypeNilValue",
			value:    &shared.RuntimeValue{Type: shared.Class, Value: nil},
			expected: true,
		},

		// 12. Default/Unknown Type (defaults to truthy as per implementation comment)
		{
			name:     "UnknownType",
			value:    &shared.RuntimeValue{Type: shared.ValueType(999), Value: "some data"},
			expected: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsTruthy(tc.value)
			if got != tc.expected {
				var valTypeStr string
				var valStr string
				if tc.value != nil {
					valTypeStr = shared.Stringify(tc.value.Type)
					valStr = fmt.Sprintf("%+v", tc.value.Value)
					t.Errorf("IsTruthy(%s: %s) = %v; want %v", valTypeStr, valStr, got, tc.expected)
				} else {
					t.Errorf("IsTruthy(nil) = %v; want %v", got, tc.expected)
				}
			}
		})
	}
}
