package lexer_test

import (
	"reflect"
	"testing"

	"github.com/dev-kas/virtlang-go/v3/lexer"
)

func TestTokenFactory(t *testing.T) {
	// "hello" (length 5) starts at L1,C1. Ends after C5. EndCol is C6.
	tok := lexer.NewToken("hello", lexer.Comment, 1, 1, 1, 6) // startLine, startCol, endLine, endCol

	if tok.Type != lexer.Comment {
		t.Errorf("expected type %v, got %v", lexer.Comment, tok.Type)
	}
	if tok.Literal != "hello" {
		t.Errorf("expected literal 'hello', got %v", tok.Literal)
	}
	if tok.StartLine != 1 || tok.StartCol != 1 || tok.EndLine != 1 || tok.EndCol != 6 {
		t.Errorf("expected L1C1-L1C6, got L%dC%d-L%dC%d", tok.StartLine, tok.StartCol, tok.EndLine, tok.EndCol)
	}
}

func TestUnescapeString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Simple Unescape",
			input:   `"hello"`,
			want:    `"hello"`,
			wantErr: false,
		},
		{
			name:    "Direct Unescape",
			input:   `hello`,
			want:    `hello`,
			wantErr: false,
		},
		{
			name:    "Escape Sequences",
			input:   `"\n\t"`,
			want:    "\"\n\t\"",
			wantErr: false,
		},
		{
			name:    "Empty String",
			input:   `""`,
			want:    `""`,
			wantErr: false,
		},
		{
			name:    "Invalid String",
			input:   `"hello`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "More Escape Sequences",
			input:   `\n\t\r`,
			want:    "\n\t\r",
			wantErr: false,
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			got, err := lexer.UnescapeString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnescapeString() %d error = %v, wantErr %v", i+1, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnescapeString() %d got = %v, want %v", i+1, got, tt.want)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []lexer.Token
		wantErr bool
	}{
		{
			name:  "String Escape Sequences",
			input: `"\n\t\r"`,
			want: []lexer.Token{
				lexer.NewToken("\n\t\r", lexer.String, 1, 1, 1, 9),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 9, 1, 9),
			},
			wantErr: false,
		},
		{
			name:  "Multiple Escape Sequences",
			input: `"Hello\nWorld\t!\"'"`,
			want: []lexer.Token{
				lexer.NewToken("Hello\nWorld\t!\"'", lexer.String, 1, 1, 1, 21),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 21, 1, 21),
			},
			wantErr: false,
		},
		{
			name:  "All Escape Sequences",
			input: `"\\\\n\t\r\b\f\\\"'"`,
			want: []lexer.Token{
				lexer.NewToken("\\\\n\t\r\b\f\\\"'", lexer.String, 1, 1, 1, 21),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 21, 1, 21),
			},
			wantErr: false,
		},
		{
			name:  "Unicode Escape Sequence",
			input: `"\u2388 <- UNICODE"`,
			want: []lexer.Token{
				lexer.NewToken("\u2388 <- UNICODE", lexer.String, 1, 1, 1, 20),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 20, 1, 20),
			},
			wantErr: false,
		},
		{
			name:  "Simple Arithmetic Tokenization",
			input: "(4+2)*3",
			want: []lexer.Token{
				lexer.NewToken("(", lexer.OParen, 1, 1, 1, 2),
				lexer.NewToken("4", lexer.Number, 1, 2, 1, 3),
				lexer.NewToken("+", lexer.BinOperator, 1, 3, 1, 4),
				lexer.NewToken("2", lexer.Number, 1, 4, 1, 5),
				lexer.NewToken(")", lexer.CParen, 1, 5, 1, 6),
				lexer.NewToken("*", lexer.BinOperator, 1, 6, 1, 7),
				lexer.NewToken("3", lexer.Number, 1, 7, 1, 8),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 8, 1, 8),
			},
			wantErr: false,
		},
		{
			name:  "Simple Lexical Tokenization",
			input: "let abc = 'hi'", // "let abc = 'hi'"
			want: []lexer.Token{ // L  C  L  C
				lexer.NewToken("let", lexer.Let, 1, 1, 1, 4),        // "let"
				lexer.NewToken("abc", lexer.Identifier, 1, 5, 1, 8), // " abc"
				lexer.NewToken("=", lexer.Equals, 1, 9, 1, 10),      // " = "
				lexer.NewToken("hi", lexer.String, 1, 11, 1, 15),  // " 'hi'"
				lexer.NewToken("<EOF>", lexer.EOF, 1, 15, 1, 15),
			},
			wantErr: false,
		},

		{
			name:  "Comprehensive Operators and Numbers",
			input: "a = (5. * .2) + 1; // Test\n b != -1.5",
			want: []lexer.Token{
				lexer.NewToken("a", lexer.Identifier, 1, 1, 1, 2),
				lexer.NewToken("=", lexer.Equals, 1, 3, 1, 4),
				lexer.NewToken("(", lexer.OParen, 1, 5, 1, 6),
				lexer.NewToken("5", lexer.Number, 1, 6, 1, 7),
				lexer.NewToken(".", lexer.Dot, 1, 7, 1, 8),
				lexer.NewToken("*", lexer.BinOperator, 1, 9, 1, 10),
				lexer.NewToken(".2", lexer.Number, 1, 11, 1, 13),
				lexer.NewToken(")", lexer.CParen, 1, 13, 1, 14),
				lexer.NewToken("+", lexer.BinOperator, 1, 15, 1, 16),
				lexer.NewToken("1", lexer.Number, 1, 17, 1, 18),
				lexer.NewToken(";", lexer.SemiColon, 1, 18, 1, 19),
				lexer.NewToken("b", lexer.Identifier, 2, 2, 2, 3),
				lexer.NewToken("!=", lexer.ComOperator, 2, 4, 2, 6),
				lexer.NewToken("-", lexer.BinOperator, 2, 7, 2, 8),
				lexer.NewToken("1.5", lexer.Number, 2, 8, 2, 11),
				lexer.NewToken("<EOF>", lexer.EOF, 2, 11, 2, 11),
			},
			wantErr: false,
		},
		{
			name:  "Multi-line Comment and Whitespace",
			input: "  let x = 10; /* Multi\n line \n comment */ const y = 20;",
			want: []lexer.Token{
				lexer.NewToken("let", lexer.Let, 1, 3, 1, 6),
				lexer.NewToken("x", lexer.Identifier, 1, 7, 1, 8),
				lexer.NewToken("=", lexer.Equals, 1, 9, 1, 10),
				lexer.NewToken("10", lexer.Number, 1, 11, 1, 13),
				lexer.NewToken(";", lexer.SemiColon, 1, 13, 1, 14),
				lexer.NewToken("const", lexer.Const, 3, 13, 3, 18),
				lexer.NewToken("y", lexer.Identifier, 3, 19, 3, 20),
				lexer.NewToken("=", lexer.Equals, 3, 21, 3, 22),
				lexer.NewToken("20", lexer.Number, 3, 23, 3, 25),
				lexer.NewToken(";", lexer.SemiColon, 3, 25, 3, 26),
				lexer.NewToken("<EOF>", lexer.EOF, 3, 26, 3, 26),
			},
			wantErr: false,
		},
		{
			name:  "String Variations",
			input: `a = "double"; b = ''; c = "with space";`,
			want: []lexer.Token{
				lexer.NewToken("a", lexer.Identifier, 1, 1, 1, 2),
				lexer.NewToken("=", lexer.Equals, 1, 3, 1, 4),
				lexer.NewToken(`double`, lexer.String, 1, 5, 1, 13),
				lexer.NewToken(";", lexer.SemiColon, 1, 13, 1, 14),
				lexer.NewToken("b", lexer.Identifier, 1, 15, 1, 16),
				lexer.NewToken("=", lexer.Equals, 1, 17, 1, 18),
				lexer.NewToken("", lexer.String, 1, 19, 1, 21),
				lexer.NewToken(";", lexer.SemiColon, 1, 21, 1, 22),
				lexer.NewToken("c", lexer.Identifier, 1, 23, 1, 24),
				lexer.NewToken("=", lexer.Equals, 1, 25, 1, 26),
				lexer.NewToken(`with space`, lexer.String, 1, 27, 1, 39),
				lexer.NewToken(";", lexer.SemiColon, 1, 39, 1, 40),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 40, 1, 40),
			},
			wantErr: false,
		},
		{
			name:  "Identifier Variations",
			input: `let _a = $b1; try c$ = 1_000;`,
			want: []lexer.Token{
				lexer.NewToken("let", lexer.Let, 1, 1, 1, 4),
				lexer.NewToken("_a", lexer.Identifier, 1, 5, 1, 7),
				lexer.NewToken("=", lexer.Equals, 1, 8, 1, 9),
				lexer.NewToken("$b1", lexer.Identifier, 1, 10, 1, 13),
				lexer.NewToken(";", lexer.SemiColon, 1, 13, 1, 14),
				lexer.NewToken("try", lexer.Try, 1, 15, 1, 18),
				lexer.NewToken("c$", lexer.Identifier, 1, 19, 1, 21),
				lexer.NewToken("=", lexer.Equals, 1, 22, 1, 23),
				lexer.NewToken("1", lexer.Number, 1, 24, 1, 25),
				lexer.NewToken("_000", lexer.Identifier, 1, 25, 1, 29),
				lexer.NewToken(";", lexer.SemiColon, 1, 29, 1, 30),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 30, 1, 30),
			},
			wantErr: false,
		},
		{
			name:  "Edge_Case_Numbers",
			input: "x = 0; y = .5; z = 5.;",
			want: []lexer.Token{
				lexer.NewToken("x", lexer.Identifier, 1, 1, 1, 2),
				lexer.NewToken("=", lexer.Equals, 1, 3, 1, 4),
				lexer.NewToken("0", lexer.Number, 1, 5, 1, 6),
				lexer.NewToken(";", lexer.SemiColon, 1, 6, 1, 7),
				lexer.NewToken("y", lexer.Identifier, 1, 8, 1, 9),
				lexer.NewToken("=", lexer.Equals, 1, 10, 1, 11),
				lexer.NewToken(".5", lexer.Number, 1, 12, 1, 14),
				lexer.NewToken(";", lexer.SemiColon, 1, 14, 1, 15),
				lexer.NewToken("z", lexer.Identifier, 1, 16, 1, 17),
				lexer.NewToken("=", lexer.Equals, 1, 18, 1, 19),
				lexer.NewToken("5", lexer.Number, 1, 20, 1, 21),
				lexer.NewToken(".", lexer.Dot, 1, 21, 1, 22),
				lexer.NewToken(";", lexer.SemiColon, 1, 22, 1, 23),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 23, 1, 23),
			},
			wantErr: false,
		},
		{
			name:  "Empty Input",
			input: "",
			want: []lexer.Token{
				lexer.NewToken("<EOF>", lexer.EOF, 1, 1, 1, 1),
			},
			wantErr: false,
		},
		{
			name:  "Whitespace Only Input",
			input: "  \n\t \r ",
			want: []lexer.Token{
				lexer.NewToken("<EOF>", lexer.EOF, 3, 2, 3, 2),
			},
			wantErr: false,
		},
		{
			name:    "Error Unrecognized Character",
			input:   "let a = #1;",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Error Unclosed String",
			input:   `a = "hello`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Error Unclosed Multi-line Comment",
			input:   `let x = /* start`,
			want:    nil,
			wantErr: true,
		},
		{
			name:  "Floating Point Number with Leading Zero",
			input: "x = 0.123",
			want: []lexer.Token{
				lexer.NewToken("x", lexer.Identifier, 1, 1, 1, 2),
				lexer.NewToken("=", lexer.Equals, 1, 3, 1, 4),
				lexer.NewToken("0.123", lexer.Number, 1, 5, 1, 10),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 10, 1, 10),
			},
			wantErr: false,
		},
		{
			name:  "Floating Point Number without Leading Zero",
			input: "y = .456",
			want: []lexer.Token{
				lexer.NewToken("y", lexer.Identifier, 1, 1, 1, 2),
				lexer.NewToken("=", lexer.Equals, 1, 3, 1, 4),
				lexer.NewToken(".456", lexer.Number, 1, 5, 1, 9),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 9, 1, 9),
			},
			wantErr: false,
		},
		{
			name:  "Floating Point Number with Trailing Dot",
			input: "z = 789.",
			want: []lexer.Token{
				lexer.NewToken("z", lexer.Identifier, 1, 1, 1, 2),
				lexer.NewToken("=", lexer.Equals, 1, 3, 1, 4),
				lexer.NewToken("789", lexer.Number, 1, 5, 1, 8),
				lexer.NewToken(".", lexer.Dot, 1, 8, 1, 9),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 9, 1, 9),
			},
			wantErr: false,
		},
		{
			name:  "Floating Point Number in Expression",
			input: "a = 1.2 + 3.4",
			want: []lexer.Token{
				lexer.NewToken("a", lexer.Identifier, 1, 1, 1, 2),
				lexer.NewToken("=", lexer.Equals, 1, 3, 1, 4),
				lexer.NewToken("1.2", lexer.Number, 1, 5, 1, 8),
				lexer.NewToken("+", lexer.BinOperator, 1, 9, 1, 10),
				lexer.NewToken("3.4", lexer.Number, 1, 11, 1, 14),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 14, 1, 14),
			},
			wantErr: false,
		},
		{
			name:    "Floating Point Number with Multiple Dots Error",
			input:   "let b = 1.2.3",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "String with newline",
			input: "let s = \"line1\\nline2\";",
			want: []lexer.Token{
				lexer.NewToken("let", lexer.Let, 1, 1, 1, 4),
				lexer.NewToken("s", lexer.Identifier, 1, 5, 1, 6),
				lexer.NewToken("=", lexer.Equals, 1, 7, 1, 8),
				// Literal is "\"line1\\nline2\"" (14 chars)
				// Content "line1\\nline2" (12 chars: l,i,n,e,1, BACKSLASH, n, l,i,n,e,2)
				// Start L1C9. Consume ". curL=1, curC=10.
				// "line1" (5 chars) -> curL=1, curC=15
				// "\" (1 char) -> curL=1, curC=16
				// "n" (1 char) -> curL=1, curC=17
				// "line2" (5 chars) -> curL=1, curC=22
				// Consume ". curL=1, curC=23.
				// So, EndLine=1, EndCol=23.
				lexer.NewToken("line1\nline2", lexer.String, 1, 9, 1, 23),
				lexer.NewToken(";", lexer.SemiColon, 1, 23, 1, 24),
				lexer.NewToken("<EOF>", lexer.EOF, 1, 24, 1, 24),
			},
			wantErr: false,
		},
		{
			name: "String with CR LF",
			// In Go, '\r\n' in a raw string literal directly inserts CR and LF runes.
			// Lexer should process CR, then LF, updating line/col correctly.
			// The token literal should contain both \r and \n.
			input: "let t = 'ab\r\ncd';",
			want: []lexer.Token{
				lexer.NewToken("let", lexer.Let, 1, 1, 1, 4),
				lexer.NewToken("t", lexer.Identifier, 1, 5, 1, 6),
				lexer.NewToken("=", lexer.Equals, 1, 7, 1, 8),
				// Literal should be "'ab\r\ncd'"
				// Content 'ab\r\ncd' (6 chars)
				// ' Start L1C9. Consume '. curL=1, curC=10.
				// a -> curC=11
				// b -> curC=12
				// \r -> curL=2, curC=1 (also consumes \n that follows)
				// c -> curL=2, curC=2
				// d -> curL=2, curC=3
				// ' Consume '. curL=2, curC=4
				// So, EndLine=2, EndCol=4
				lexer.NewToken("ab\ncd", lexer.String, 1, 9, 2, 4),
				lexer.NewToken(";", lexer.SemiColon, 2, 4, 2, 5),
				lexer.NewToken("<EOF>", lexer.EOF, 2, 5, 2, 5),
			},
			wantErr: false,
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := lexer.Tokenize(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Tokenize() %d error = %v, wantErr %v", i+1, err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					if len(got) != len(tt.want) {
						t.Errorf("Tokenize() %d token count mismatch:\nexpected: %d tokens\n     got: %d tokens", i+1, len(tt.want), len(got))
						t.Logf("Expected tokens: %+v", tt.want)
						t.Logf("Got tokens:      %+v", got)
						for i := 0; i < min(len(got), len(tt.want)); i++ {
							if !reflect.DeepEqual(got[i], tt.want[i]) {
								t.Errorf("Tokenize() %d mismatch (first diff):\nexpected: %+v\n     got: %+v", i+1, tt.want[i], got[i])
							}
						}
						if len(got) > len(tt.want) {
							t.Errorf("Extra tokens from got: %+v", got[len(tt.want):])
						} else if len(tt.want) > len(got) {
							t.Errorf("Missing tokens, expected: %+v", tt.want[len(got):])
						}
					} else {
						for j := range got {
							if !reflect.DeepEqual(got[j], tt.want[j]) {
								t.Errorf("Tokenize() %d mismatch at index %d:\nexpected: %+v (%s)\n     got: %+v (%s)",
									i+1, j, tt.want[j], lexer.Stringify(tt.want[j].Type), got[j], lexer.Stringify(got[j].Type))
								break
							}
						}
						// Fallback, though the loop above should catch it.
						t.Errorf("Tokenize() general mismatch (check individual elements if not caught above):\nexpected: %+v\n     got: %+v", tt.want, got)
					}
				}
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
