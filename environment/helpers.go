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
	var newEnv *Environment
	if env.Parent != nil {
		parentCopy := DeepCopy(env.Parent)
		envCopy := NewEnvironment(parentCopy)
		envCopy.Global = env.Global
		newEnv = &envCopy
	} else {
		envCopy := NewEnvironment(nil)
		envCopy.Global = env.Global
		newEnv = &envCopy
	}

	// Copy variables
	for k, v := range env.Variables {
		// Create a deep copy of the RuntimeValue
		var newValue any
		switch val := v.Value.(type) {
		case map[string]*shared.RuntimeValue: // Object
			newMap := make(map[string]*shared.RuntimeValue, len(val))
			for mk, mv := range val {
				if mv != nil {
					mvCopy := *mv
					newMap[mk] = &mvCopy
				}
			}
			newValue = newMap
		case []shared.RuntimeValue: // Array
			newSlice := make([]shared.RuntimeValue, len(val))
			copy(newSlice, val)
			newValue = newSlice
		default:
			// For primitive types, we can use the value as is
			newValue = v.Value
		}

		(*newEnv).Variables[k] = shared.RuntimeValue{
			Type:  v.Type,
			Value: newValue,
		}
	}

	// Copy constants
	for k, v := range env.Constants {
		(*newEnv).Constants[k] = v
	}

	return newEnv
}
