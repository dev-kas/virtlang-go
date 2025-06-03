package debugger_test

import (
	"sync"
	"testing"
	"time"

	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/debugger"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/shared"
)

// TestNewDebugger tests the creation of a new debugger
func TestNewDebugger(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(env)

	if dbg == nil {
		t.Fatal("NewDebugger returned nil")
	}

	// Check initial state
	if dbg.State != debugger.RunningState {
		t.Errorf("Expected initial state to be RunningState, got %v", dbg.State)
	}

	// Check environment is set correctly
	if dbg.Environment != env {
		t.Error("Environment not set correctly")
	}
}

// TestStateTransitions tests the state transition methods
func TestStateTransitions(t *testing.T) {
	tests := []struct {
		name      string
		operation func(*debugger.Debugger) error
		setup     func(*debugger.Debugger)
		expected  debugger.State
	}{
		{"Pause", (*debugger.Debugger).Pause, nil, debugger.PausedState},
		{"StepInto", (*debugger.Debugger).StepInto, nil, debugger.SteppingState},
		{"Continue", (*debugger.Debugger).Continue, func(d *debugger.Debugger) {
			d.Pause()
		}, debugger.RunningState},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			env := environment.NewEnvironment(nil)
			dbg := debugger.NewDebugger(env)

			if tc.setup != nil {
				tc.setup(dbg)
			}

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
	dbg := debugger.NewDebugger(env)

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
	dbg := debugger.NewDebugger(env)

	// The debugger should use the provided environment
	if dbg.Environment != env {
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

// TestWaitIfPaused tests the WaitIfPaused functionality
func TestWaitIfPaused(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(env)

	// Test that WaitIfPaused doesn't block when not paused
	done := make(chan bool)
	go func() {
		dbg.WaitIfPaused(ast.FnDeclarationNode)
		done <- true
	}()

	select {
	case <-done:
		// Expected - should not block
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitIfPaused blocked when not paused")
	}

	// Test that WaitIfPaused blocks when paused
	dbg.Pause()
	go func() {
		dbg.WaitIfPaused(ast.FnDeclarationNode)
		done <- true
	}()

	select {
	case <-done:
		t.Fatal("WaitIfPaused didn't block when paused")
	case <-time.After(100 * time.Millisecond):
		// Expected - should block
	}

	// Test that Continue unblocks WaitIfPaused
	dbg.Continue()

	select {
	case <-done:
		// Expected - should unblock
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitIfPaused didn't unblock after Continue")
	}
}

// TestIsDebuggable tests the IsDebuggable method
func TestIsDebuggable(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(env)

	tests := []struct {
		name     string
		nodeType ast.NodeType
		expected bool
	}{
		{"VarDeclaration", ast.VarDeclarationNode, true},
		{"FnDeclaration", ast.FnDeclarationNode, true},
		{"IfStatement", ast.IfStatementNode, true},
		{"WhileLoop", ast.WhileLoopNode, true},
		{"ReturnStmt", ast.ReturnStmtNode, true},
		{"ContinueStmt", ast.ContinueStmtNode, true},
		{"BreakStmt", ast.BreakStmtNode, true},
		{"TryCatchStmt", ast.TryCatchStmtNode, true},
		{"CallExpr", ast.CallExprNode, true},
		{"Class", ast.ClassNode, true},
		{"ClassMethod", ast.ClassMethodNode, true},
		{"ClassProperty", ast.ClassPropertyNode, true},
		{"Program", ast.ProgramNode, true},
		{"Unknown", 9999, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := dbg.IsDebuggable(tc.nodeType)
			if result != tc.expected {
				t.Errorf("Expected IsDebuggable(%s) to be %v, got %v", tc.name, tc.expected, result)
			}
		})
	}
}

// TestConcurrentAccess tests thread safety of the debugger
func TestConcurrentAccess(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(env)

	var wg sync.WaitGroup
	numRoutines := 10

	// Test concurrent state changes
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dbg.Pause()
			dbg.Continue()
			dbg.StepInto()
		}()
	}

	// Test concurrent WaitIfPaused
	wg.Add(1)
	go func() {
		defer wg.Done()
		dbg.WaitIfPaused(ast.FnDeclarationNode)
	}()

	// Let the goroutines run for a bit
	time.Sleep(100 * time.Millisecond)

	// Unblock any waiting goroutines
	dbg.Continue()

	wg.Wait()
}

// TestStepInto tests the StepInto functionality with correct state expectations
func TestStepInto(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(env)

	// Set up a goroutine to simulate execution
	done := make(chan bool)
	go func() {
		dbg.WaitIfPaused(ast.FnDeclarationNode)
		done <- true
	}()

	// Start stepping into
	if err := dbg.StepInto(); err != nil {
		t.Fatalf("StepInto failed: %v", err)
	}

	// Verify state is SteppingState after StepInto
	if dbg.State != debugger.SteppingState {
		t.Errorf("Expected state SteppingState, got %v", dbg.State)
	}

	// Continue execution
	if err := dbg.Continue(); err != nil {
		t.Fatalf("Continue failed: %v", err)
	}

	// Wait for the goroutine to finish
	select {
	case <-done:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Test timed out")
	}
}

// TestStepOver tests the StepOver functionality with correct state expectations
func TestStepOver(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(env)

	// Set up initial call stack
	dbg.PushFrame(debugger.StackFrame{Filename: "main.go", Line: 1})

	// Set up a goroutine to simulate execution
	done := make(chan bool)
	go func() {
		dbg.WaitIfPaused(ast.FnDeclarationNode)
		done <- true
	}()

	// Start stepping over
	if err := dbg.StepOver(); err != nil {
		t.Fatalf("StepOver failed: %v", err)
	}

	// Verify state is SteppingState after StepOver
	if dbg.State != debugger.SteppingState {
		t.Errorf("Expected state SteppingState, got %v", dbg.State)
	}

	// Continue execution
	if err := dbg.Continue(); err != nil {
		t.Fatalf("Continue failed: %v", err)
	}

	// Wait for the goroutine to finish
	select {
	case <-done:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Test timed out")
	}
}

// TestStepOut tests the StepOut functionality with correct state expectations
func TestStepOut(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(env)

	// Set up initial call stack with multiple frames
	dbg.PushFrame(debugger.StackFrame{Filename: "main.go", Line: 1})
	dbg.PushFrame(debugger.StackFrame{Filename: "utils.go", Line: 5})

	// Set up a goroutine to simulate execution
	done := make(chan bool)
	go func() {
		dbg.WaitIfPaused(ast.FnDeclarationNode)
		done <- true
	}()

	// Start stepping out
	if err := dbg.StepOut(); err != nil {
		t.Fatalf("StepOut failed: %v", err)
	}

	// Verify state is SteppingState after StepOut
	if dbg.State != debugger.SteppingState {
		t.Errorf("Expected state SteppingState, got %v", dbg.State)
	}

	// Continue execution
	if err := dbg.Continue(); err != nil {
		t.Fatalf("Continue failed: %v", err)
	}

	// Wait for the goroutine to finish
	select {
	case <-done:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Test timed out")
	}
}

// TestCallStack tests the call stack functionality
func TestCallStack(t *testing.T) {
	env := environment.NewEnvironment(nil)
	dbg := debugger.NewDebugger(env)

	// Test initial call stack is empty
	if len(dbg.CallStack) != 0 {
		t.Fatalf("Expected empty call stack, got length %d", len(dbg.CallStack))
	}

	// Push a frame
	frame1 := debugger.StackFrame{Filename: "test1.go", Line: 10}
	dbg.PushFrame(frame1)

	if len(dbg.CallStack) != 1 {
		t.Fatalf("Expected call stack length 1, got %d", len(dbg.CallStack))
	}

	if dbg.CallStack[0].Filename != "test1.go" || dbg.CallStack[0].Line != 10 {
		t.Fatalf("Unexpected frame data: %+v", dbg.CallStack[0])
	}

	// Push another frame
	frame2 := debugger.StackFrame{Filename: "test2.go", Line: 20}
	dbg.PushFrame(frame2)

	if len(dbg.CallStack) != 2 {
		t.Fatalf("Expected call stack length 2, got %d", len(dbg.CallStack))
	}

	// Pop a frame
	dbg.PopFrame()
	if len(dbg.CallStack) != 1 {
		t.Fatalf("Expected call stack length 1 after pop, got %d", len(dbg.CallStack))
	}

	// Pop last frame
	dbg.PopFrame()
	if len(dbg.CallStack) != 0 {
		t.Fatalf("Expected empty call stack after pops, got length %d", len(dbg.CallStack))
	}

	// Test pop empty stack (should not panic)
	dbg.PopFrame()
}
