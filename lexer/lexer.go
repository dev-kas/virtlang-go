package lexer

import (
	"strings"
	"unicode"

	"github.com/dev-kas/VirtLang-Go/errors"
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
	ComOperator                  // < == > != <= =>
	If                           // if
	Else                         // else
	String                       // '...' "..."
	WhileLoop                    // while
	Comment                      // --<...>-- -->...
	Try                          // try
	Catch                        // catch
	Return                       // return
	Break                        // break
	Continue                     // continue
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
}

type Token struct {
	Type       TokenType
	Literal    string
	Start      int
	Difference int
}

func NewToken(value string, tokenType TokenType, start, diff int) Token {
	return Token{
		Type:       tokenType,
		Literal:    value,
		Start:      start,
		Difference: diff,
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
	src := strings.Split(srcCode, "")
	srcLen := len(src)

	position := 0

	for position < srcLen {
		start := position
		char := src[position]
		currentCharRune := rune(char[0])

		if IsSkippable(currentCharRune) {
			position++
			continue
		}

		tokenLen := 1
		nextToken := Token{}
		createToken := true

		switch char {
		case "(":
			nextToken = NewToken(char, OParen, start+1, tokenLen)
		case ")":
			nextToken = NewToken(char, CParen, start+1, tokenLen)
		case "{":
			nextToken = NewToken(char, OBrace, start+1, tokenLen)
		case "}":
			nextToken = NewToken(char, CBrace, start+1, tokenLen)
		case "[":
			nextToken = NewToken(char, OBracket, start+1, tokenLen)
		case "]":
			nextToken = NewToken(char, CBracket, start+1, tokenLen)
		case ";":
			nextToken = NewToken(char, SemiColon, start+1, tokenLen)
		case ":":
			nextToken = NewToken(char, Colon, start+1, tokenLen)
		case ",":
			nextToken = NewToken(char, Comma, start+1, tokenLen)

		default:
			createToken = false
		}

		if createToken {
			tokens = append(tokens, nextToken)
			position += tokenLen
			continue
		}

		if char == "/" && position+1 < srcLen {
			nextChar := src[position+1]
			if nextChar == "/" {
				position += 2
				for position < srcLen && src[position] != "\n" && src[position] != "\r" {
					position++
				}
				continue
			} else if nextChar == "*" {
				position += 2
				foundEnd := false
				for position+1 < srcLen {
					if src[position] == "*" && src[position+1] == "/" {
						foundEnd = true
						position += 2
						break
					}
					position++
				}
				if !foundEnd {
					return nil, &errors.LexerError{Character: ' ', Position: start + 1}
				}
				continue
			}
		}

		if IsBinaryOperator(currentCharRune) {
			tokens = append(tokens, NewToken(char, BinOperator, start+1, 1))
			position++
			continue
		}

		if strings.ContainsRune("=<>!", currentCharRune) {
			if position+1 < srcLen {
				twoCharOp := char + src[position+1]
				if IsComparisonOperator(twoCharOp) {
					tokens = append(tokens, NewToken(twoCharOp, ComOperator, start+1, 2))
					position += 2
					continue
				}
			}
			if IsComparisonOperator(char) {
				tokens = append(tokens, NewToken(char, ComOperator, start+1, 1))
				position++
				continue
			} else if char == "=" {
				tokens = append(tokens, NewToken(char, Equals, start+1, 1))
				position++
				continue
			}
		}

		if currentCharRune == '\'' || currentCharRune == '"' {
			quote := currentCharRune
			position++
			strContent := ""
			stringStart := start + 1
			foundEndQuote := false
			for position < srcLen {
				if rune(src[position][0]) == quote {
					foundEndQuote = true
					break
				}
				// TODO: Handle escape sequences
				strContent += src[position]
				position++
			}

			if foundEndQuote {
				position++
				fullLiteral := string(quote) + strContent + string(quote)
				tokens = append(tokens, NewToken(fullLiteral, String, stringStart, position-start))
			} else {
				return nil, &errors.LexerError{Character: quote, Position: stringStart}
			}
			continue
		}

		if currentCharRune == '.' {
			if position+1 < srcLen && IsNumeric(rune(src[position+1][0])) {
				numStr := "."
				position++
				numStartPos := start
				for position < srcLen && IsNumeric(rune(src[position][0])) {
					numStr += src[position]
					position++
				}
				tokens = append(tokens, NewToken(numStr, Number, numStartPos+1, position-numStartPos))
			} else {

				tokens = append(tokens, NewToken(char, Dot, start+1, 1))
				position++
			}
			continue
		}

		if IsNumeric(currentCharRune) {
			numStart := start
			numStr := ""
			hasDecimal := false
			for position < srcLen {
				loopChar := src[position]
				loopRune := rune(loopChar[0])

				if IsNumeric(loopRune) {
					numStr += loopChar
					position++
				} else if loopChar == "." {
					if hasDecimal {
						return nil, &errors.LexerError{
							Character: '.',
							Position:  position + 1,
						}
					}
					if position+1 < srcLen && IsNumeric(rune(src[position+1][0])) {
						numStr += loopChar
						hasDecimal = true
						position++
					} else {
						break
					}
				} else {
					break
				}
			}
			tokens = append(tokens, NewToken(numStr, Number, numStart+1, position-numStart))
			continue
		}

		if IsAlpha(currentCharRune) || currentCharRune == '_' || currentCharRune == '$' {
			identStart := start
			identStr := ""

			for position < srcLen {
				loopChar := src[position]
				loopRune := rune(loopChar[0])
				if IsAlphaNumeric(loopRune) || loopRune == '_' || loopRune == '$' {
					identStr += loopChar
					position++
				} else {
					break
				}
			}

			if tokenType, isKeyword := KEYWORDS[identStr]; isKeyword {
				tokens = append(tokens, NewToken(identStr, tokenType, identStart+1, position-identStart))
			} else {
				tokens = append(tokens, NewToken(identStr, Identifier, identStart+1, position-identStart))
			}
			continue
		}

		return nil, &errors.LexerError{
			Character: currentCharRune,
			Position:  start + 1,
		}

	}

	tokens = append(tokens, NewToken("<EOF>", EOF, position+1, 0))

	return tokens, nil
}
