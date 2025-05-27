package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/dev-kas/virtlang-go/v3/errors"
)

type TokenType int

const (
	Number      TokenType = iota // 0 - 9
	Identifier                   // a - z A - Z 0 - 9 _ $
	Equals                       // =
	BinOperator                  // / * + _
	OParen                       // (
	CParen                       // )
	Let                          // let
	Const                        // const
	SemiColon                    // ;
	Comma                        // ,
	Colon                        // :
	OBrace                       // {
	CBrace                       // }
	OBracket                     // [
	CBracket                     // ]
	Dot                          // .
	Fn                           // fn
	ComOperator                  // < == > != <= >=
	If                           // if
	Else                         // else
	String                       // '...' "..."
	WhileLoop                    // while
	Comment                      // /*...*/ //...
	Try                          // try
	Catch                        // catch
	Return                       // return
	Break                        // break
	Continue                     // continue
	Class                        // class
	Public                       // public
	Private                      // private
	EOF                          // end of file
)

func Stringify(t TokenType) string {
	switch t {
	case Number:
		return "Number"
	case Identifier:
		return "Identifier"
	case Equals:
		return "Equals"
	case BinOperator:
		return "BinOperator"
	case OParen:
		return "OParen"
	case CParen:
		return "CParen"
	case Let:
		return "Let"
	case Const:
		return "Const"
	case SemiColon:
		return "SemiColon"
	case Comma:
		return "Comma"
	case Colon:
		return "Colon"
	case OBrace:
		return "OBrace"
	case CBrace:
		return "CBrace"
	case OBracket:
		return "OBracket"
	case CBracket:
		return "CBracket"
	case Dot:
		return "Dot"
	case Fn:
		return "Fn"
	case ComOperator:
		return "ComOperator"
	case If:
		return "If"
	case Else:
		return "Else"
	case String:
		return "String"
	case WhileLoop:
		return "WhileLoop"
	case Comment:
		return "Comment"
	case Try:
		return "Try"
	case Catch:
		return "Catch"
	case Return:
		return "Return"
	case Break:
		return "Break"
	case Continue:
		return "Continue"
	case Class:
		return "Class"
	case Public:
		return "Public"
	case Private:
		return "Private"
	case EOF:
		return "EOF"
	default:
		return "Unknown"
	}
}

var KEYWORDS = map[string]TokenType{
	"let":      Let,
	"const":    Const,
	"fn":       Fn,
	"if":       If,
	"else":     Else,
	"while":    WhileLoop,
	"try":      Try,
	"catch":    Catch,
	"return":   Return,
	"break":    Break,
	"continue": Continue,
	"class":    Class,
	"public":   Public,
	"private":  Private,
}

type Token struct {
	Type      TokenType
	Literal   string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
}

func NewToken(value string, tokenType TokenType, startLine, startCol, endLine, endCol int) Token {
	return Token{
		Type:      tokenType,
		Literal:   value,
		StartLine: startLine,
		StartCol:  startCol,
		EndLine:   endLine,
		EndCol:    endCol,
	}
}

func IsAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

func IsNumeric(r rune) bool {
	return unicode.IsNumber(r)
}

func IsAlphaNumeric(r rune) bool {
	return IsAlpha(r) || IsNumeric(r)
}

func IsSkippable(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func IsBinaryOperator(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '%':
		return true
	default:
		return false
	}
}

func UnescapeString(s string) (string, error) {
	stringQuoted := s
	isQuoted := false
	
	if len(stringQuoted) == 0 {
		return "", nil
	}

	if stringQuoted[0] == '"' && stringQuoted[len(stringQuoted)-1] == '"' {
		isQuoted = true
	} else {
		stringQuoted = fmt.Sprintf("\"%s\"", stringQuoted)
	}

	unquoted, err := strconv.Unquote(stringQuoted)
	if err != nil {
		return "", err
	}

	stringQuoted = unquoted

	if !isQuoted {
		return stringQuoted, nil
	}

	return fmt.Sprintf("\"%s\"", stringQuoted), nil
}

func IsComparisonOperator(r string) bool {
	switch r {
	case "==", ">", "<", "!=", "<=", ">=":
		return true
	default:
		return false
	}
}

func Tokenize(srcCode string) ([]Token, *errors.LexerError) {
	var tokens []Token
	runes := []rune(srcCode) // Process source as a slice of runes
	srcLen := len(runes)

	position := 0      // Current index in `runes`
	currentLine := 1   // 1-based line number
	currentColumn := 1 // 1-based column number

	// Use rune as key for single character tokens
	singleCharTokenMap := map[rune]TokenType{
		'(': OParen,
		')': CParen,
		'{': OBrace,
		'}': CBrace,
		'[': OBracket,
		']': CBracket,
		';': SemiColon,
		':': Colon,
		',': Comma,
	}

	for position < srcLen {
		tokStartLine := currentLine
		tokStartCol := currentColumn

		currentCharRune := runes[position] // Get the current rune

		// --- 1. Skippable characters ---
		if IsSkippable(currentCharRune) {
			position++ // Consume the skippable rune

			if currentCharRune == '\n' {
				currentLine++
				currentColumn = 1
			} else if currentCharRune == '\r' {
				currentLine++
				currentColumn = 1
				// Handle \r\n sequence
				if position < srcLen && runes[position] == '\n' {
					position++ // Consume \n as part of the same line break
				}
			} else { // space or tab
				currentColumn++
			}
			continue
		}

		// --- 2. Single-character punctuation ---
		if tokenType, ok := singleCharTokenMap[currentCharRune]; ok {
			position++      // Consume the rune
			currentColumn++ // These characters don't cause line breaks
			tokens = append(tokens, NewToken(string(currentCharRune), tokenType, tokStartLine, tokStartCol, currentLine, currentColumn))
			continue
		}

		// --- 3. Comments (// and /* */) ---
		if currentCharRune == '/' {
			if position+1 < srcLen {
				nextRune := runes[position+1]
				if nextRune == '/' { // Single-line comment: //
					position += 2 // Consume "//"
					currentColumn += 2

					for position < srcLen {
						consumedCommentRune := runes[position]
						position++ // Consume character in comment body

						if consumedCommentRune == '\n' {
							currentLine++
							currentColumn = 1
							break
						} else if consumedCommentRune == '\r' {
							currentLine++
							currentColumn = 1
							if position < srcLen && runes[position] == '\n' { // Check for \r\n
								position++
							}
							break
						} else {
							currentColumn++
						}
					}
					continue
				} else if nextRune == '*' { // Multi-line comment: /* ... */
					position += 2 // Consume "/*"
					currentColumn += 2
					unclosedCommentErrorPos := errors.Position{Line: tokStartLine, Col: tokStartCol}
					foundEnd := false
					for position+1 < srcLen { // Need at least two runes for "*/"
						char1Rune := runes[position]

						position++ // Consume current rune in comment body
						if char1Rune == '\n' {
							currentLine++
							currentColumn = 1
						} else if char1Rune == '\r' {
							currentLine++
							currentColumn = 1
							if position < srcLen && runes[position] == '\n' { // Check for \r\n
								position++
							}
						} else {
							currentColumn++
						}

						if char1Rune == '*' && runes[position] == '/' {
							foundEnd = true
							position++ // Consume "/"
							currentColumn++
							break
						}
					}
					if !foundEnd {
						return nil, &errors.LexerError{
							Character: '/',
							Pos:       unclosedCommentErrorPos,
						}
					}
					continue
				}
			}
		}

		// --- 4. Binary Operators (standalone +, -, *, /, %) ---
		if IsBinaryOperator(currentCharRune) {
			position++
			currentColumn++
			tokens = append(tokens, NewToken(string(currentCharRune), BinOperator, tokStartLine, tokStartCol, currentLine, currentColumn))
			continue
		}

		// --- 5. Comparison and Equals Operators (=, ==, >, <, !=, <=, >=) ---
		if strings.ContainsRune("=<>!", currentCharRune) {
			firstOpCharStr := string(currentCharRune)

			if position+1 < srcLen {
				secondOpRune := runes[position+1]
				twoCharOp := firstOpCharStr + string(secondOpRune)
				if IsComparisonOperator(twoCharOp) {
					position += 2
					currentColumn += 2
					tokens = append(tokens, NewToken(twoCharOp, ComOperator, tokStartLine, tokStartCol, currentLine, currentColumn))
					continue
				}
			}

			if IsComparisonOperator(firstOpCharStr) { // <, >
				position++
				currentColumn++
				tokens = append(tokens, NewToken(firstOpCharStr, ComOperator, tokStartLine, tokStartCol, currentLine, currentColumn))
				continue
			} else if firstOpCharStr == "=" { // =
				position++
				currentColumn++
				tokens = append(tokens, NewToken(firstOpCharStr, Equals, tokStartLine, tokStartCol, currentLine, currentColumn))
				continue
			}
		}

		// --- 6. String Literals ('...' or "...") ---
		if currentCharRune == '\'' || currentCharRune == '"' {
			quoteRune := currentCharRune
			openingQuoteStr := string(currentCharRune)

			position++      // Consume opening quote from runes
			currentColumn++ // Opening quote advances column

			unclosedStringErrorPos := errors.Position{Line: tokStartLine, Col: tokStartCol}
			var strContentBuilder strings.Builder // Builds the content BETWEEN quotes
			foundEndQuote := false

			for position < srcLen {
				loopRune := runes[position] // Current rune from source

				if loopRune == quoteRune && runes[position-1] != '\\' { // Not an escaped quote
					position++      // Consume the closing quote from runes
					currentColumn++ // The closing quote itself advances column
					foundEndQuote = true
					break
				}

				strContentBuilder.WriteRune(loopRune) // Add the current rune to our string's content
				position++                            // Consume the rune from source

				if loopRune == '\n' { // LF
					currentLine++
					currentColumn = 1
				} else if loopRune == '\r' { // CR
					currentLine++
					currentColumn = 1
					if position < srcLen && runes[position] == '\n' { // Check for CRLF
						strContentBuilder.WriteRune(runes[position]) // Add LF part of CRLF
						position++                                   // Consume the LF from runes
					}
				} else { // Regular character
					currentColumn++
				}
			}

			if foundEndQuote {
				stringified := strContentBuilder.String()
				finalBuilder := strings.Builder{}
				splitted := strings.Split(strings.ReplaceAll(strings.ReplaceAll(stringified, "\r\n", "\n"), "\r", "\n"), "\n")
				for i, line := range splitted {
					unescapedLiteral, err := UnescapeString(line)
					if err != nil {
						return nil, &errors.LexerError{
							Character: quoteRune,
							Pos:       unclosedStringErrorPos,
							Message:   err.Error(),
						}
					}
					finalBuilder.WriteString(unescapedLiteral)
					if i != len(splitted)-1 {
						finalBuilder.WriteString("\n")
					}
				}
				fullLiteral := openingQuoteStr + finalBuilder.String() + string(quoteRune)
				tokens = append(tokens, NewToken(fullLiteral, String, tokStartLine, tokStartCol, currentLine, currentColumn))
			} else {
				return nil, &errors.LexerError{
					Character: quoteRune,
					Pos:       unclosedStringErrorPos,
				}
			}
			continue
		}

		// --- 7. Numbers (integers, floats, including dot-prefixed like .5) ---
		if currentCharRune == '.' {
			if position+1 < srcLen && IsNumeric(runes[position+1]) { // Number like .5
				numStartIndex := position
				position++ // Consume '.'
				currentColumn++

				for position < srcLen && IsNumeric(runes[position]) {
					position++
					currentColumn++
				}
				literal := string(runes[numStartIndex:position])
				tokens = append(tokens, NewToken(literal, Number, tokStartLine, tokStartCol, currentLine, currentColumn))
				continue
			} else { // Just a Dot token
				position++
				currentColumn++
				tokens = append(tokens, NewToken(".", Dot, tokStartLine, tokStartCol, currentLine, currentColumn))
				continue
			}
		}

		if IsNumeric(currentCharRune) {
			numStartIndex := position
			hasDecimal := false
			for position < srcLen {
				loopRune := runes[position]

				if IsNumeric(loopRune) {
					position++
					currentColumn++
				} else if loopRune == '.' {
					if hasDecimal {
						return nil, &errors.LexerError{
							Character: '.',
							Pos:       errors.Position{Line: currentLine, Col: currentColumn},
						}
					}
					if position+1 < srcLen && IsNumeric(runes[position+1]) {
						hasDecimal = true
						position++ // Consume '.'
						currentColumn++
					} else { // Number ends before '.', e.g., "123."
						break
					}
				} else { // Not a digit or a valid part of a number
					break
				}
			}
			literal := string(runes[numStartIndex:position])
			tokens = append(tokens, NewToken(literal, Number, tokStartLine, tokStartCol, currentLine, currentColumn))
			continue
		}

		// --- 8. Identifiers and Keywords ---
		if IsAlpha(currentCharRune) || currentCharRune == '_' || currentCharRune == '$' {
			identStartIndex := position

			for position < srcLen {
				loopRune := runes[position]
				if IsAlphaNumeric(loopRune) || loopRune == '_' || loopRune == '$' {
					position++
					currentColumn++
				} else {
					break // End of identifier/keyword
				}
			}
			literal := string(runes[identStartIndex:position])
			tokenType, isKeyword := KEYWORDS[literal]
			if !isKeyword {
				tokenType = Identifier
			}
			tokens = append(tokens, NewToken(literal, tokenType, tokStartLine, tokStartCol, currentLine, currentColumn))
			continue
		}

		// --- 9. Unrecognized Character ---
		return nil, &errors.LexerError{
			Character: currentCharRune,
			Pos:       errors.Position{Line: tokStartLine, Col: tokStartCol},
		}
	} // End of main for loop

	// --- EOF Token ---
	tokens = append(tokens, NewToken("<EOF>", EOF, currentLine, currentColumn, currentLine, currentColumn))
	return tokens, nil
}
