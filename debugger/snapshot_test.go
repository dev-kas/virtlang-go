package debugger_test

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v3/debugger"
	"github.com/dev-kas/virtlang-go/v3/environment"
	"github.com/dev-kas/virtlang-go/v3/shared"
)

func TestDebugger_TakeSnapshot(t *testing.T) {
	t.Run("takes a snapshot of current stack and environment", func(t *testing.T) {
		env := environment.NewEnvironment(nil)
		env.DeclareVar("testVar", shared.RuntimeValue{Type: shared.Number, Value: 42.0}, false)

		dbg := debugger.NewDebugger(&env)

		dbg.PushFrame(debugger.StackFrame{
			Name:     "testFunction",
			Filename: "test.virt",
			Line:     10,
		})

		dbg.TakeSnapshot()

		if len(dbg.Snapshots) != 1 {
			t.Fatalf("expected 1 snapshot, got %d", len(dbg.Snapshots))
		}

		snapshot := dbg.Snapshots[0]
		if len(snapshot.Stack) != 1 {
			t.Fatalf("expected 1 frame in snapshot, got %d", len(snapshot.Stack))
		}

		frame := snapshot.Stack[0]
		if frame.Name != "testFunction" {
			t.Errorf("expected function name 'testFunction', got '%s'", frame.Name)
		}
		if frame.Filename != "test.virt" {
			t.Errorf("expected filename 'test.virt', got '%s'", frame.Filename)
		}

		// Verify the environment has the same content
		val, err := snapshot.Env.LookupVar("testVar")
		if err != nil {
			t.Fatalf("failed to lookup testVar in snapshot: %v", err)
		}
		if val.Value != 42.0 {
			t.Errorf("expected testVar to be 42.0 in snapshot, got %v", val.Value)
		}
	})

	t.Run("snapshot remains unchanged after environment or stack modifications", func(t *testing.T) {
		// Create initial environment with a variable
		env := environment.NewEnvironment(nil)
		env.DeclareVar("testVar", shared.RuntimeValue{Type: shared.Number, Value: 42.0}, false)

		dbg := debugger.NewDebugger(&env)

		// Push initial frame
		dbg.PushFrame(debugger.StackFrame{
			Name:     "testFunction",
			Filename: "test.virt",
			Line:     10,
		})

		// Take snapshot of current state
		dbg.TakeSnapshot()

		// Mutate the existing environment
		// Update the existing variable using AssignVar
		_, err := env.AssignVar("testVar", shared.RuntimeValue{Type: shared.Number, Value: 100.0})
		if err != nil {
			t.Fatalf("failed to update testVar: %v", err)
		}
		// Add a new variable
		env.DeclareVar("newVar", shared.RuntimeValue{Type: shared.String, Value: "modified"}, false)

		// Push a new frame to test stack depth
		dbg.PushFrame(debugger.StackFrame{
			Name:     "newFunction",
			Filename: "new.virt",
			Line:     5,
		})

		// Modify the top frame
		dbg.CallStack[0].Line = 20
		dbg.CallStack[0].Name = "modifiedFunction"

		// Verify snapshot remains unchanged
		snapshot := dbg.Snapshots[0]

		// Check environment in snapshot is unchanged
		snapshotVal, err := snapshot.Env.LookupVar("testVar")
		if err != nil {
			t.Fatalf("failed to lookup testVar in snapshot: %v", err)
		}
		if snapshotVal.Value != 42.0 {
			t.Errorf("expected testVar to be 42.0 in snapshot, got %v", snapshotVal.Value)
		}

		_, err = snapshot.Env.LookupVar("newVar")
		if err == nil {
			t.Error("expected newVar to not exist in snapshot")
		}

		// Check stack frame in snapshot is unchanged
		if len(snapshot.Stack) != 1 {
			t.Fatalf("expected 1 frame in snapshot, got %d", len(snapshot.Stack))
		}

		frame := snapshot.Stack[0]
		if frame.Name != "testFunction" {
			t.Errorf("expected function name 'testFunction' in snapshot, got '%s'", frame.Name)
		}
		if frame.Line != 10 {
			t.Errorf("expected line 10 in snapshot, got %d", frame.Line)
		}

		// Verify current environment reflects changes
		currentVal, err := env.LookupVar("testVar")
		if err != nil {
			t.Fatalf("failed to lookup testVar in current env: %v", err)
		}
		if currentVal.Value != 100.0 {
			t.Errorf("expected testVar to be 100.0 in current env, got %v", currentVal.Value)
		}

		// Verify current stack reflects changes
		if len(dbg.CallStack) > 0 && dbg.CallStack[0].Name != "modifiedFunction" {
			t.Errorf("expected current function name to be 'modifiedFunction', got '%s'", dbg.CallStack[0].Name)
		}

		// Verify the top frame was modified
		topFrame := dbg.CallStack[0]
		expectedName := "modifiedFunction"
		expectedLine := 20
		if topFrame.Name != expectedName || topFrame.Line != expectedLine {
			t.Errorf("expected top frame to be %s:%d, got %s:%d",
				expectedName, expectedLine, topFrame.Name, topFrame.Line)
		}

		// Verify current environment has the modified values
		currentVal, err2 := dbg.Environment.LookupVar("testVar")
		if err2 != nil || currentVal.Value != 100.0 {
			t.Errorf("expected testVar to be 100.0 in current environment, got %v", currentVal)
		}

		// Verify first frame was modified as expected
		if dbg.CallStack[0].Name != expectedName || dbg.CallStack[0].Line != expectedLine {
			t.Errorf("expected first stack frame to be %s:%d, got %v",
				expectedName, expectedLine, dbg.CallStack[0])
		}

		// Verify stack frame mutations
		if len(dbg.CallStack) < 1 || dbg.CallStack[0].Name != "modifiedFunction" || dbg.CallStack[0].Line != 20 {
			t.Errorf("expected first stack frame to be modifiedFunction:20, got %+v", dbg.CallStack[0])
		}
	})
}
