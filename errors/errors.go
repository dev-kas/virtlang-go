package errors

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v3/shared"
)

type Position struct {
	Line int // 1-based line number
	Col  int // 1-based column number
}

// --- InternalCommunicationProtocol ---
type InternalCommunicationProtocolTypes int

const (
	ICP_Return InternalCommunicationProtocolTypes = iota
	ICP_Continue
	ICP_Break
)

type InternalCommunicationProtocol struct {
	Type   InternalCommunicationProtocolTypes
	RValue *shared.RuntimeValue
}

// --- RuntimeError ---
type RuntimeError struct {
	Message                       string
	InternalCommunicationProtocol *InternalCommunicationProtocol
}

func (e *RuntimeError) Error() string {
	if e.Message == "" {
		e.Message = "Unspecified error" // Or perhaps "Unknown runtime error"
	}
	return fmt.Sprintf("Runtime Error: %s", e.Message)
}

// --- SyntaxError ---
type SyntaxError struct {
	Expected string   // What was expected
	Got      string   // What was actually found (e.g., token literal or type)
	Start    Position // Start position of the problematic syntax
	End      Position // End position of the problematic syntax
	Message  string   // Optional additional message
}

func (e *SyntaxError) Error() string {
	var positionStr string
	// Check if it's a point error or a very short range on the same line
	if e.Start.Line == e.End.Line && e.End.Col <= e.Start.Col+1 { // e.g. single char/token, or empty range
		positionStr = fmt.Sprintf("at L%dC%d", e.Start.Line, e.Start.Col)
	} else {
		positionStr = fmt.Sprintf("from L%dC%d to L%dC%d", e.Start.Line, e.Start.Col, e.End.Line, e.End.Col)
	}

	baseMsg := "Syntax Error"
	if e.Message != "" {
		baseMsg = fmt.Sprintf("%s: %s", baseMsg, e.Message)
	}

	if e.Expected != "" && e.Got != "" {
		return fmt.Sprintf("%s: Expected %s, got '%s' %s", baseMsg, e.Expected, e.Got, positionStr)
	} else if e.Got != "" { // Only 'Got' is specified
		return fmt.Sprintf("%s: Unexpected '%s' %s", baseMsg, e.Got, positionStr)
	} else { // Neither Expected nor Got, or only Expected (less common for "got")
		return fmt.Sprintf("%s %s", baseMsg, positionStr)
	}
}

// NewSyntaxError creates a new SyntaxError.
// 'gotLiteral' is often the token.Literal that caused the error.
// 'start' and 'end' define the span of the problematic token/syntax.
func NewSyntaxError(expected string, gotLiteral string, start, end Position) *SyntaxError {
	return &SyntaxError{
		Expected: expected,
		Got:      gotLiteral,
		Start:    start,
		End:      end,
	}
}

// NewSyntaxErrorf creates a new SyntaxError with a custom message.
func NewSyntaxErrorf(start, end Position, format string, args ...interface{}) *SyntaxError {
	return &SyntaxError{
		Start:   start,
		End:     end,
		Message: fmt.Sprintf(format, args...),
	}
}

// --- ParserError ---
// Often, ParserError is very similar to SyntaxError.
type ParserError struct {
	Token   string   // The literal of the unexpected token
	Start   Position // Start position of the token
	End     Position // End position of the token
	Message string   // Optional: More specific message about why it's an error
}

func (e *ParserError) Error() string {
	var positionStr string
	if e.Start.Line == e.End.Line && e.End.Col <= e.Start.Col+1 {
		positionStr = fmt.Sprintf("at L%dC%d", e.Start.Line, e.Start.Col)
	} else {
		positionStr = fmt.Sprintf("from L%dC%d to L%dC%d", e.Start.Line, e.Start.Col, e.End.Line, e.End.Col)
	}

	baseMsg := "Parser Error"
	if e.Message != "" {
		baseMsg = fmt.Sprintf("%s: %s", baseMsg, e.Message)
	}

	if e.Token != "" {
		return fmt.Sprintf("%s: Unexpected token '%s' %s", baseMsg, e.Token, positionStr)
	}
	return fmt.Sprintf("%s %s", baseMsg, positionStr)
}

// NewParserError creates a new ParserError for an unexpected token.
func NewParserError(tokenLiteral string, start, end Position) *ParserError {
	return &ParserError{
		Token: tokenLiteral,
		Start: start,
		End:   end,
	}
}

// NewParserErrorf creates a new ParserError with a custom formatted message.
func NewParserErrorf(start, end Position, format string, args ...interface{}) *ParserError {
	return &ParserError{
		Start:   start,
		End:     end,
		Message: fmt.Sprintf(format, args...),
	}
}

// --- LexerError ---
type LexerError struct {
	Character rune     // The problematic character
	Pos       Position // Position (Line/Col) of the character
	Message   string   // Optional: for specific messages like "unclosed comment"
}

func (e *LexerError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Lexer Error: %s at L%dC%d", e.Message, e.Pos.Line, e.Pos.Col)
	}

	return fmt.Sprintf("Lexer Error: Unexpected character '%c' at L%dC%d", e.Character, e.Pos.Line, e.Pos.Col)
}

// NewLexerError creates a new LexerError for an unexpected character (Message will be empty).
func NewLexerError(char rune, pos Position) *LexerError {
	return &LexerError{Character: char, Pos: pos}
}

// NewLexerErrorf creates a new LexerError with a custom message.
// 'charForContext' can be the opening delimiter (e.g., '/' for unclosed comment) or 0 if not relevant.
// The Message field will be populated by the formatted string.
func NewLexerErrorf(pos Position, charForContext rune, format string, args ...interface{}) *LexerError {
	return &LexerError{
		Pos:       pos,
		Character: charForContext,
		Message:   fmt.Sprintf(format, args...),
	}
}
