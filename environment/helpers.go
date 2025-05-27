package environment

import (
	"github.com/dev-kas/virtlang-go/v3/shared"
)

// DeepCopy creates a deep copy of the environment
func DeepCopy(env *Environment) *Environment {
	if env == nil {
		return nil
	}

	// Create a new environment with the same parent
	newEnv := &Environment{
		Parent:    DeepCopy(env.Parent),
		Variables: make(map[string]*shared.RuntimeValue),
		Constants: make(map[string]struct{}),
		Global:    env.Global,
	}

	// Copy variables
	for k, v := range env.Variables {
		if v == nil {
			continue
		}

		var newValue interface{}
		switch val := v.Value.(type) {
		case map[string]*shared.RuntimeValue:
			// Deep copy the object
			newMap := make(map[string]*shared.RuntimeValue)
			for k, v := range val {
				if v != nil {
					valCopy := *v
					newMap[k] = &valCopy
				}
			}
			newValue = newMap
		case []shared.RuntimeValue:
			// Deep copy the array
			newSlice := make([]shared.RuntimeValue, len(val))
			copy(newSlice, val)
			newValue = newSlice
		default:
			// For primitive types, we can use the value as is
			newValue = v.Value
		}

		// Create a new RuntimeValue pointer
		rv := &shared.RuntimeValue{
			Type:  v.Type,
			Value: newValue,
		}
		(*newEnv).Variables[k] = rv
	}

	// Copy constants
	for k, v := range env.Constants {
		(*newEnv).Constants[k] = v
	}

	return newEnv
}
