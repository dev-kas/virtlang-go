package debugger_test

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/debugger"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/shared"
)

// TestNewDebugger tests the creation of a new debugger
func TestNewDebugger(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(&env)

	if dbg == nil {
		t.Fatal("NewDebugger returned nil")
	}

	// Check initial state
	if dbg.State != debugger.RunningState {
		t.Errorf("Expected initial state to be RunningState, got %v", dbg.State)
	}

	// Check environment is set correctly
	if dbg.Environment != &env {
		t.Error("Environment not set correctly")
	}
}

// TestStateTransitions tests the state transition methods
func TestStateTransitions(t *testing.T) {
	tests := []struct {
		name      string
		operation func(*debugger.Debugger) error
		expected  debugger.State
	}{
		{"Run", (*debugger.Debugger).Run, debugger.RunningState},
		{"Pause", (*debugger.Debugger).Pause, debugger.PausedState},
		{"Step", (*debugger.Debugger).Step, debugger.SteppingState},
		{"Continue", (*debugger.Debugger).Continue, debugger.RunningState},
		{"Stop", (*debugger.Debugger).Stop, debugger.RunningState},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			env := environment.NewEnvironment(nil)
			dbg := debugger.NewDebugger(&env)

			err := tc.operation(dbg)
			if err != nil {
				t.Fatalf("%s returned error: %v", tc.name, err)
			}

			if dbg.State != tc.expected {
				t.Errorf("Expected state %v after %s, got %v", tc.expected, tc.name, dbg.State)
			}
		})
	}
}

// TestBreakpointIntegration tests integration with BreakpointManager
func TestBreakpointIntegration(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(&env)

	file := "test.go"
	line := 42

	// Initially no breakpoint
	if dbg.ShouldStop(file, line) {
		t.Fatal("Should not stop at line without breakpoint")
	}

	// Add breakpoint
	dbg.BreakpointManager.Set(file, line)

	// Should stop at breakpoint
	if !dbg.ShouldStop(file, line) {
		t.Fatal("Should stop at line with breakpoint")
	}

	// Different line should not stop
	if dbg.ShouldStop(file, line+1) {
		t.Fatal("Should not stop at different line")
	}
}

// TestEnvironmentIntegration tests that the environment is passed through correctly
func TestEnvironmentIntegration(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(&env)

	// The debugger should use the provided environment
	if dbg.Environment != &env {
		t.Error("Debugger does not use the provided environment")
	}

	// Changes to the environment should be visible through the debugger
	_, err := env.DeclareVar("test", shared.RuntimeValue{
		Type:  shared.Number,
		Value: 42.0,
	}, true)
	if err != nil {
		t.Fatalf("Failed to declare variable: %v", err)
	}

	val, err := dbg.Environment.LookupVar("test")
	if err != nil {
		t.Fatalf("Failed to lookup variable: %v", err)
	}

	if val.Type != shared.Number || val.Value != 42.0 {
		t.Errorf("Environment value not as expected, got type %v, value %v", val.Type, val.Value)
	}
}

// TestInterfaceCompleteness tests that all required methods are implemented
func TestInterfaceCompleteness(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(&env)

	// This test will fail to compile if the Debugger doesn't implement all required methods
	var _ interface {
		Run() error
		Pause() error
		Stop() error
		Step() error
		StepOut() error
		StepInto() error
		Continue() error
		ShouldStop(string, int) bool
	} = dbg
}

// TestIsDebuggable tests the IsDebuggable method
func TestIsDebuggable(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(&env)

	// Test debuggable node types
	debuggableNodes := []ast.Stmt{
		&ast.VarDeclaration{},
		&ast.VarAssignmentExpr{},
		&ast.IfStatement{},
		&ast.WhileLoop{},
		&ast.ReturnStmt{},
		&ast.ContinueStmt{},
		&ast.BreakStmt{},
		&ast.TryCatchStmt{},
		&ast.CallExpr{},
		&ast.FnDeclaration{},
		&ast.Class{},
		&ast.ClassMethod{},
		&ast.ClassProperty{},
		&ast.Program{},
	}

	for _, node := range debuggableNodes {
		nodeType := node.GetType().String()
		t.Run("should be debuggable: "+nodeType, func(t *testing.T) {
			if !dbg.IsDebuggable(node) {
				t.Errorf("Expected node type %s to be debuggable", nodeType)
			}
		})
	}

	// Test non-debuggable node types (just a few examples)
	nonDebuggableNodes := []ast.Stmt{
		&ast.BinaryExpr{},
		&ast.CompareExpr{},
		&ast.Identifier{},
		&ast.NumericLiteral{},
		&ast.StringLiteral{},
	}

	for _, node := range nonDebuggableNodes {
		nodeType := node.GetType().String()
		t.Run("should not be debuggable: "+nodeType, func(t *testing.T) {
			if dbg.IsDebuggable(node) {
				t.Errorf("Expected node type %s to not be debuggable", nodeType)
			}
		})
	}
}
