package environment_test

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/shared"
)

func TestNewEnvironment(t *testing.T) {
	global := environment.NewEnvironment(nil)

	if global.Parent != nil {
		t.Errorf("Expected Parent to be nil, but got %v", global.Parent)
	}

	local := environment.NewEnvironment(&global)

	if local.Parent != &global {
		t.Errorf("Local environment parent should be global, but got %v and %v", &global, local.Parent)
	}
}

func TestDeclareVariable(t *testing.T) {
	env := environment.NewEnvironment(nil)
	val := shared.RuntimeValue{
		Type:  shared.Number,
		Value: 10,
	}

	declaredVar, err := env.DeclareVar("x", val, false)
	if err != nil {
		t.Errorf("Expected err to be nil, but got %v", err)
	}
	if declaredVar == nil || declaredVar.Value != val.Value {
		t.Errorf("Expected declaredVar to be %v, but got %v", val, declaredVar)
	}

	lookupVal, err := env.LookupVar("x")
	if err != nil {
		t.Errorf("Expected err to be nil, but got %v", err)
	}
	if lookupVal == nil || lookupVal.Value != val.Value {
		t.Errorf("Expected lookupVal to be %v, but got %v", val, lookupVal)
	}

	constantVal := shared.RuntimeValue{
		Type:  shared.String,
		Value: "hello",
	}
	_, err = env.DeclareVar("GREETING", constantVal, true)
	if err != nil {
		t.Errorf("Expected err to be nil, but got %v", err)
	}
	lookedUpConstant, err := env.LookupVar("GREETING")
	if err != nil {
		t.Errorf("Expected err to be nil, but got %v", err)
	}
	if lookedUpConstant == nil || lookedUpConstant.Value != constantVal.Value {
		t.Errorf("Expected lookedUpConstant to be %v, but got %v", constantVal, lookedUpConstant)
	}

	// declare variable already declared
	_, err = env.DeclareVar(
		"x",
		shared.RuntimeValue{
			Type:  shared.Number,
			Value: 20,
		},
		false,
	)
	if err == nil {
		t.Error("Expected err, but got nil")
	}
}

func TestAssignVar(t *testing.T) {
	env := environment.NewEnvironment(nil)
	env.DeclareVar("y", shared.RuntimeValue{Type: shared.Number, Value: 5}, false)

	newValue := shared.RuntimeValue{Type: shared.Number, Value: 15}
	assignedVal, err := env.AssignVar("y", newValue)
	if err != nil {
		t.Errorf("AssignVar returned error: %v", err)
	}
	if *assignedVal != newValue {
		t.Errorf("AssignVar returned incorrect value. Expected: %v, Got: %v", newValue, assignedVal)
	}

	lookupVal, err := env.LookupVar("y")
	if err != nil {
		t.Errorf("LookupVar after AssignVar returned error: %v", err)
	}
	if lookupVal == nil || lookupVal.Value != newValue.Value {
		t.Errorf("LookupVar after AssignVar returned incorrect value. Expected: %v, Got: %v", newValue, lookupVal)
	}

	// Test assignment to constant panics
	_, err = env.DeclareVar("PI", shared.RuntimeValue{Type: shared.Number, Value: 3}, true)
	if err != nil {
		t.Errorf("AssignVar to constant should have panicked")
	}
	env.AssignVar("PI", shared.RuntimeValue{Type: shared.Number, Value: 314})
}

func TestResolveAndLookupVarWithParent(t *testing.T) {
	global := environment.NewEnvironment(nil)
	global.DeclareVar("globalVar", shared.RuntimeValue{Type: shared.String, Value: "global value"}, false)

	local := environment.NewEnvironment(&global)
	local.DeclareVar("localVar", shared.RuntimeValue{Type: shared.String, Value: "local value"}, false)

	// Lookup in local scope
	localLookup, err := local.LookupVar("localVar")
	if err != nil {
		t.Errorf("LookupVar in local scope returned error: %v", err)
	}
	if localLookup == nil || localLookup.Value != "local value" {
		t.Errorf("LookupVar in local scope failed. Expected: %v, Got: %v", "local value", localLookup)
	}

	// Lookup in parent (global) scope
	globalLookupFromLocal, err := local.LookupVar("globalVar")
	if err != nil {
		t.Errorf("LookupVar in parent scope returned error: %v", err)
	}
	if globalLookupFromLocal == nil || globalLookupFromLocal.Value != "global value" {
		t.Errorf("LookupVar in parent scope failed. Expected: %v, Got: %v", "global value", globalLookupFromLocal)
	}

	// Test resolving
	resolvedLocal, err := local.Resolve("localVar")
	if err != nil {
		t.Errorf("resolve(\"localVar\") returned error: %v", err)
	}
	if resolvedLocal != &local {
		t.Errorf("resolve(\"localVar\") should have returned the local environment")
	}

	resolvedGlobal, err := local.Resolve("globalVar")
	if err != nil {
		t.Errorf("resolve(\"globalVar\") returned error: %v", err)
	}
	if resolvedGlobal != &global {
		t.Errorf("resolve(\"globalVar\") should have returned the global environment")
	}

	// Test resolving non-existent variable panics
	_, err = local.Resolve("nonExistentVar")
	if err == nil {
		t.Errorf("resolve(\"nonExistentVar\") should have panicked")
	}

	// Test looking up non-existent variable panics
	_, err = local.LookupVar("nonExistentVar")
	if err == nil {
		t.Errorf("LookupVar(\"nonExistentVar\") should have panicked")
	}
}
