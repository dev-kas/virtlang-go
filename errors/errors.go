package errors

import "fmt"

// RuntimeError
type RuntimeError struct {
	Message string
}

func (e *RuntimeError) Error() string {
	if e.Message == "" {
		e.Message = "Unspecified error"
	}

	return fmt.Sprintf(
		"Runtime Error: %s",
		e.Message,
	)
}

// SyntaxError
type SyntaxError struct {
	Expected   string
	Got        string
	Start      int
	Difference int
}

func (e *SyntaxError) Error() string {
	errFormatStart := "Syntax Error: Expected %s got %s"
	errFormatEnd := "from position %d to position %d"

	mode := 0

	if e.Expected == "" || e.Got == "" {
		errFormatStart = "Syntax Error: Unspecified error"
		mode++
	}

	if e.Difference <= 1 {
		errFormatEnd = "at position %d"
		mode += 2
	}

	if mode == 0 { // Syntax Error: Expected %s got %s from position %d to position %d"
		return fmt.Sprintf(
			errFormatStart+" "+errFormatEnd,
			e.Expected,
			e.Got,
			e.Start,
			e.Difference+e.Start,
		)
	} else if mode == 1 { // Syntax Error: Unspecified error from position %d to position %d
		return fmt.Sprintf(
			errFormatStart+" "+errFormatEnd,
			e.Start,
			e.Difference+e.Start,
		)
	} else if mode == 2 { // Syntax Error: Expected %s got %s at position %d
		return fmt.Sprintf(
			errFormatStart+" "+errFormatEnd,
			e.Expected,
			e.Got,
			e.Start,
		)
	} else { // Syntax Error: Unspecified error at position %d
		return fmt.Sprintf(
			errFormatStart+" "+errFormatEnd,
			e.Start,
		)
	}
}

// ParserError
type ParserError struct {
	Token      string
	Start      int
	Difference int
}

func (e *ParserError) Error() string {
	errFormatStart := "Parser Error: Unexpected token %s"
	errFormatEnd := "from position %d to position %d"

	mode := 0

	if e.Token == "" {
		errFormatStart = "Parser Error: Unspecified error"
		mode++
	}

	if e.Difference <= 1 {
		errFormatEnd = "at position %d"
		mode += 2
	}

	if mode == 0 { // Parser Error: Unexpected token %s from position %d to position %d
		return fmt.Sprintf(
			errFormatStart+" "+errFormatEnd,
			e.Token,
			e.Start,
			e.Difference+e.Start,
		)
	} else if mode == 1 { // Parser Error: Unspecified error from position %d to position %d
		return fmt.Sprintf(
			errFormatStart+" "+errFormatEnd,
			e.Start,
			e.Difference+e.Start,
		)
	} else if mode == 2 { // Parser Error: Unexpected token %s at position %d
		return fmt.Sprintf(
			errFormatStart+" "+errFormatEnd,
			e.Token,
			e.Start,
		)
	} else { // Parser Error: Unspecified error at position %d
		return fmt.Sprintf(
			errFormatStart+" "+errFormatEnd,
			e.Start,
		)
	}
}

// LexerError
type LexerError struct {
	Character rune
	Position  int
}

func (e *LexerError) Error() string {
	return fmt.Sprintf(
		"Lexer Error: Unexpected character '%c' (0x%02x) at position %d",
		e.Character,
		e.Character,
		e.Position,
	)
}
