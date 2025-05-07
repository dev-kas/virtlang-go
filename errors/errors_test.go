package errors_test

import (
	"testing"

	"github.com/dev-kas/virtlang-go/errors"
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
			want:    "Runtime Error: Unspecified error",
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
				t.Errorf("Got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSyntaxError(t *testing.T) {
	tests := []struct {
		name       string
		expected   string
		got        string
		start      int
		difference int
		want       string
	}{
		{
			name:       "Simple Syntax Error",
			expected:   "identifier",
			got:        "nil",
			start:      5,
			difference: 3,
			want:       "Syntax Error: Expected identifier got nil from position 5 to position 8",
		},
		{
			name:       "Empty Syntax Error",
			expected:   "",
			got:        "",
			start:      0,
			difference: 0,
			want:       "Syntax Error: Unspecified error at position 0",
		},
		{
			name:       "Missing Closing Brace",
			expected:   "'}'",
			got:        "EOF",
			start:      42,
			difference: 1,
			want:       "Syntax Error: Expected '}' got EOF at position 42",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := &errors.SyntaxError{
				Expected:   tt.expected,
				Got:        tt.got,
				Start:      tt.start,
				Difference: tt.difference,
			}

			if got := err.Error(); got != tt.want {
				t.Errorf("Got %q\nWant %q", got, tt.want)
			}
		})
	}
}

func TestParserError(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		start      int
		difference int
		want       string
	}{
		{
			name:       "Simple Parser Error",
			token:      "identifier",
			start:      5,
			difference: 3,
			want:       "Parser Error: Unexpected token identifier from position 5 to position 8",
		},
		{
			name:       "Empty Parser Error",
			token:      "",
			start:      0,
			difference: 0,
			want:       "Parser Error: Unspecified error at position 0",
		},
		{
			name:       "Closing Brace",
			token:      "}",
			start:      42,
			difference: 1,
			want:       "Parser Error: Unexpected token } at position 42",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := &errors.ParserError{
				Token:      tt.token,
				Start:      tt.start,
				Difference: tt.difference,
			}

			if got := err.Error(); got != tt.want {
				t.Errorf("Got %q\nWant %q", got, tt.want)
			}
		})
	}
}

func TestLexerError(t *testing.T) {
	tests := []struct {
		name string
		char rune
		pos  int
		want string
	}{
		{
			name: "Simple Lexer Error",
			char: 'a',
			pos:  5,
			want: "Lexer Error: Unexpected character 'a' (0x61) at position 5",
		},
		{
			name: "Null Unicode Lexer Error",
			char: '\000',
			pos:  0,
			want: "Lexer Error: Unexpected character '\x00' (0x00) at position 0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := &errors.LexerError{
				Character: tt.char,
				Position:  tt.pos,
			}

			if got := err.Error(); got != tt.want {
				t.Errorf("Got %q\nWant %q", got, tt.want)
			}
		})
	}
}
