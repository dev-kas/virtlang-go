package debugger_test

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v3/debugger"
)

// TestBreakpointManager tests the public API of BreakpointManager
func TestBreakpointManager(t *testing.T) {
	t.Run("should initialize with no breakpoints", func(t *testing.T) {
		bm := debugger.NewBreakpointManager()
		if bm == nil {
			t.Fatal("BreakpointManager should not be nil")
		}

		// Test that no breakpoints exist initially
		if bm.Has("any_file.go", 1) {
			t.Fatal("New BreakpointManager should not have any breakpoints")
		}
	})

	t.Run("should set and check breakpoints", func(t *testing.T) {
		bm := debugger.NewBreakpointManager()
		file := "test.go"
		line := 42

		// Initially no breakpoint
		if bm.Has(file, line) {
			t.Fatal("Breakpoint should not exist initially")
		}

		// Set breakpoint
		bm.Set(file, line)

		// Verify it exists
		if !bm.Has(file, line) {
			t.Fatal("Breakpoint should exist after being set")
		}

		// Set another breakpoint
		bm.Set(file, line+1)
		if !bm.Has(file, line+1) {
			t.Fatal("Second breakpoint should exist")
		}

		// First breakpoint should still exist
		if !bm.Has(file, line) {
			t.Fatal("First breakpoint should still exist")
		}
	})

	t.Run("should remove breakpoints", func(t *testing.T) {
		bm := debugger.NewBreakpointManager()
		file := "test.go"
		line := 42

		// Set and verify breakpoint
		bm.Set(file, line)
		if !bm.Has(file, line) {
			t.Fatal("Breakpoint should exist after being set")
		}

		// Remove breakpoint
		bm.Remove(file, line)

		// Verify it's gone
		if bm.Has(file, line) {
			t.Fatal("Breakpoint should not exist after removal")
		}

		// Removing non-existent breakpoint should not panic
		bm.Remove("nonexistent.go", 999)
	})

	t.Run("should clear all breakpoints", func(t *testing.T) {
		bm := debugger.NewBreakpointManager()

		// Add some breakpoints
		bm.Set("file1.go", 1)
		bm.Set("file2.go", 2)
		bm.Set("file3.go", 3)

		// Verify they exist
		if !bm.Has("file1.go", 1) || !bm.Has("file2.go", 2) || !bm.Has("file3.go", 3) {
			t.Fatal("All breakpoints should exist before clear")
		}

		// Clear all breakpoints
		bm.Clear()

		// Verify all are gone
		if bm.Has("file1.go", 1) || bm.Has("file2.go", 2) || bm.Has("file3.go", 3) {
			t.Fatal("No breakpoints should exist after clear")
		}
	})

	t.Run("should handle edge cases", func(t *testing.T) {
		bm := debugger.NewBreakpointManager()

		// Test empty filename
		bm.Set("", 0)
		if !bm.Has("", 0) {
			t.Fatal("Should handle empty filename")
		}

		// Test negative line numbers
		bm.Set("test.go", -1)
		if !bm.Has("test.go", -1) {
			t.Fatal("Should handle negative line numbers")
		}

		// Test removing non-existent breakpoint
		bm.Remove("nonexistent.go", 999)
	})
}
