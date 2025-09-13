package evaluator

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func evalDestructureDeclaration(node *ast.DestructureDeclaration, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	rhsValue, err := Evaluate(node.Value, env, dbgr)
	if err != nil {
		return nil, err
	}

	err = evalDestructureDeclaration_bindPattern(node.Pattern, rhsValue, env, node.Constant, dbgr)
	if err != nil {
		return nil, err
	}

	return rhsValue, nil
}

// recursive helper function
func evalDestructureDeclaration_bindPattern(pattern ast.DestructurePattern, value *shared.RuntimeValue, env *environment.Environment, isConstant bool, dbgr *debugger.Debugger) *errors.RuntimeError {
	switch p := pattern.(type) {
	case *ast.DestructureObjectPattern:
		// ensure we're destructuring an object
		if value.Type != shared.Object {
			return &errors.RuntimeError{
				Message: fmt.Sprintf("Cannot destructure value of type %s with an object pattern", shared.Stringify(value.Type)),
			}
		}

		objValue := value.Value.(map[string]*shared.RuntimeValue)
		assignedKeys := make(map[string]bool)

		for _, prop := range p.Properties {
			subValue, exists := objValue[prop.Key]

			// default values if key doesnt exist
			if !exists || subValue.Type == shared.Nil {
				if prop.Default != nil {
					evaluatedDefault, err := Evaluate(prop.Default, env, dbgr)
					if err != nil {
						return err
					}
					subValue = evaluatedDefault
				} else {
					nilVal := values.MK_NIL()
					subValue = &nilVal
				}
			}

			if prop.DeconstructChildren != nil {
				// we have a nested pattern
				err := evalDestructureDeclaration_bindPattern(prop.DeconstructChildren, subValue, env, isConstant, dbgr)
				if err != nil {
					return err
				}
			} else {
				// simple variable binding
				_, err := env.DeclareVar(prop.Name, *subValue, isConstant)
				if err != nil {
					return err
				}
			}

			assignedKeys[prop.Key] = true
		}

		// handle the rest operator (...)
		if p.Rest != nil {
			restObj := make(map[string]*shared.RuntimeValue)
			for key, val := range objValue {
				if !assignedKeys[key] {
					restObj[key] = val
				}
			}
			restValue := values.MK_OBJECT(restObj)
			_, err := env.DeclareVar(*p.Rest, restValue, isConstant)
			if err != nil {
				return err
			}
		}

		return nil

	case *ast.DestructureArrayPattern:
		var arrValue []shared.RuntimeValue 

		// check if the value is a string and convert it to an array of characters
		if value.Type == shared.String {
			// convert string to array
			str := value.Value.(string)
			charValues := make([]shared.RuntimeValue, 0, len(str))
			for _, char := range str {
				charValues = append(charValues, values.MK_STRING(string(char)))
			}
			arrValue = charValues
		} else if value.Type == shared.Array {
			// already an array
			arrValue = value.Value.([]shared.RuntimeValue)
		} else {
			// unknown
			return &errors.RuntimeError{
				Message: fmt.Sprintf("Cannot destructure value of type %s with an array pattern", shared.Stringify(value.Type)),
			}
		}

		// tterate through the pattern's elements
		for i, element := range p.Elements {
			// if the pattern has a blank element like [a, , c] then skip it
			if element.Skipped {
				continue
			}

			var subValue *shared.RuntimeValue

			// check if source array has an element at this index
			if i >= len(arrValue) {
				// out of bounds: use default value or nil
				if element.Default != nil {
					evaluatedDefault, err := Evaluate(element.Default, env, dbgr)
					if err != nil {
						return err
					}
					subValue = evaluatedDefault
				} else {
					nilVal := values.MK_NIL()
					subValue = &nilVal
				}
			} else {
				// in bounds: get the value from the source array
				subValue = &arrValue[i]
			}

			if element.DeconstructChildren != nil {
				err := evalDestructureDeclaration_bindPattern(element.DeconstructChildren, subValue, env, isConstant, dbgr)
				if err != nil {
					return err
				}
			} else {
				_, err := env.DeclareVar(element.Name, *subValue, isConstant)
				if err != nil {
					return err
				}
			}
		}

		// handle the rest operator (...) for arrays
		if p.Rest != nil {
			startIndex := len(p.Elements)
			var restSlice []shared.RuntimeValue

			if startIndex < len(arrValue) {
				restSlice = arrValue[startIndex:]
			} else {
				restSlice = []shared.RuntimeValue{} // empty slice if no elements are left
			}
			
			restValue := values.MK_ARRAY(restSlice)
			_, err := env.DeclareVar(*p.Rest, restValue, isConstant)
			if err != nil {
				return err
			}
		}

		return nil

	default:
		return &errors.RuntimeError{Message: fmt.Sprintf("Unhandled pattern type: %T", p)}
	}
}
