package errors_test

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v3/errors"
	// "github.com/dev-kas/virtlang-go/v3/shared" // We'll use errors.Position
)

func TestRuntimeError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "Simple Runtime Error",
			message: "An error occurred",
			want:    "Runtime Error: An error occurred",
		},
		{
			name:    "Empty Runtime Error",
			message: "",
			// Updated based on the provided errors.go (it defaults to "Unspecified error" if Message is empty)
			want: "Runtime Error: Unspecified error",
		},
		{
			name:    "Complex Error",
			message: "An error occurred with a lot of details",
			want:    "Runtime Error: An error occurred with a lot of details",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := &errors.RuntimeError{
				Message: tt.message,
			}

			if got := err.Error(); got != tt.want {
				t.Errorf("RuntimeError.Error()\nGot:  %q\nWant: %q", got, tt.want)
			}
		})
	}
}

func TestSyntaxError(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		got      string
		start    errors.Position // Use errors.Position
		end      errors.Position // Use errors.Position
		message  string          // For testing the message field
		want     string
	}{
		{
			name:     "Expected identifier got nil - range",
			expected: "identifier",
			got:      "nil",
			start:    errors.Position{Line: 1, Col: 5},
			end:      errors.Position{Line: 1, Col: 8}, // End.Col (8) > Start.Col+1 (6)
			want:     "Syntax Error: Expected identifier, got 'nil' from L1C5 to L1C8",
		},
		{
			name:     "Expected '}' got EOF - at specific point",
			expected: "'}'",
			got:      "EOF",
			start:    errors.Position{Line: 1, Col: 42},
			end:      errors.Position{Line: 1, Col: 43}, // End.Col (43) == Start.Col+1 (43)
			want:     "Syntax Error: Expected '}', got 'EOF' at L1C42",
		},
		{
			name:     "Expected specific token - multi-line range",
			expected: "';'",
			got:      "вар", // "var" in Cyrillic, just an example
			start:    errors.Position{Line: 2, Col: 10},
			end:      errors.Position{Line: 3, Col: 2},
			want:     "Syntax Error: Expected ';', got 'вар' from L2C10 to L3C2",
		},
		{
			name:    "Custom message error - at point",
			start:   errors.Position{Line: 5, Col: 1},
			end:     errors.Position{Line: 5, Col: 1}, // Point error
			message: "Invalid assignment target",
			want:    "Syntax Error: Invalid assignment target at L5C1",
		},
		{
			name:  "Only got specified - range",
			got:   "++",
			start: errors.Position{Line: 6, Col: 3},
			end:   errors.Position{Line: 6, Col: 5},
			want:  "Syntax Error: Unexpected '++' from L6C3 to L6C5",
		},
		{
			name:  "Neither expected nor got - just position",
			start: errors.Position{Line: 7, Col: 7},
			end:   errors.Position{Line: 7, Col: 7},
			want:  "Syntax Error at L7C7",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := &errors.SyntaxError{
				Expected: tt.expected,
				Got:      tt.got,
				Start:    tt.start,
				End:      tt.end,
				Message:  tt.message,
			}

			if got := err.Error(); got != tt.want {
				t.Errorf("SyntaxError.Error()\nGot:  %q\nWant: %q", got, tt.want)
			}
		})
	}
}

func TestParserError(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		start   errors.Position // Use errors.Position
		end     errors.Position // Use errors.Position
		message string          // For testing the message field
		want    string
	}{
		{
			name:  "Unexpected token - range",
			token: "unexpected_token_literal",
			start: errors.Position{Line: 1, Col: 10},
			end:   errors.Position{Line: 1, Col: 30}, // End.Col > Start.Col+1
			want:  "Parser Error: Unexpected token 'unexpected_token_literal' from L1C10 to L1C30",
		},
		{
			name:  "Unexpected token - at specific point",
			token: "?",
			start: errors.Position{Line: 2, Col: 5},
			end:   errors.Position{Line: 2, Col: 6}, // End.Col == Start.Col+1
			want:  "Parser Error: Unexpected token '?' at L2C5",
		},
		{
			name:    "Custom message error - range",
			start:   errors.Position{Line: 10, Col: 1},
			end:     errors.Position{Line: 10, Col: 5},
			message: "Missing semicolon after statement",
			want:    "Parser Error: Missing semicolon after statement from L10C1 to L10C5",
		},
		{
			name:  "No token, just position - at point", // Simulates an error where token info isn't primary
			start: errors.Position{Line: 11, Col: 1},
			end:   errors.Position{Line: 11, Col: 1},
			want:  "Parser Error at L11C1",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := &errors.ParserError{
				Token:   tt.token,
				Start:   tt.start,
				End:     tt.end,
				Message: tt.message,
			}

			if got := err.Error(); got != tt.want {
				t.Errorf("ParserError.Error()\nGot:  %q\nWant: %q", got, tt.want)
			}
		})
	}
}

func TestLexerError(t *testing.T) {
	tests := []struct {
		name    string
		char    rune
		pos     errors.Position
		message string
		want    string
	}{
		{
			name: "Simple Lexer Error - Unexpected char",
			char: '#',
			pos:  errors.Position{Line: 1, Col: 5},
			want: "Lexer Error: Unexpected character '#' at L1C5",
		},
		{
			name: "Null Unicode Lexer Error - Unexpected char",
			char: '\000',
			pos:  errors.Position{Line: 1, Col: 0}, // Col 0 is unusual, but for test
			want: "Lexer Error: Unexpected character '\x00' at L1C0",
		},
		{
			name:    "Lexer Error with custom message - unclosed comment",
			pos:     errors.Position{Line: 10, Col: 2},
			message: "Unclosed multi-line comment",
			char:    '/',
			want:    "Lexer Error: Unclosed multi-line comment at L10C2",
		},
		{
			name:    "Lexer Error with custom message - no specific char",
			pos:     errors.Position{Line: 11, Col: 3},
			message: "Invalid number format",
			char:    0,
			want:    "Lexer Error: Invalid number format at L11C3",
		},
		{
			name:    "Lexer Error with char 0 and no message",
			pos:     errors.Position{Line: 12, Col: 4},
			char:    0,
			message: "",
			want:    "Lexer Error: Unexpected character '\x00' at L12C4",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := &errors.LexerError{
				Character: tt.char,
				Pos:       tt.pos,
				Message:   tt.message,
			}

			if got := err.Error(); got != tt.want {
				t.Errorf("LexerError.Error()\nGot:  %q\nWant: %q", got, tt.want)
			}
		})
	}
}
