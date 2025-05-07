package lexer_test

import (
	"reflect"
	"testing"

	"github.com/dev-kas/virtlang-go/lexer"
)

func TestTokenFactory(t *testing.T) {
	tok := lexer.NewToken("hello", lexer.Comment, 1, 5)

	if tok.Type != lexer.Comment {
		t.Errorf("expected type %v, got %v", lexer.Comment, tok.Type)
	}
	if tok.Literal != "hello" {
		t.Errorf("expected literal 'hello', got %v", tok.Literal)
	}
	if tok.Start != 1 || tok.Difference != 5 {
		t.Errorf("expected start=1 and end=5, got start=%d end=%d", tok.Start, tok.Difference)
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
			name:  "Simple Arithmetic Tokenization",
			input: "(4+2)*3",
			want: []lexer.Token{
				lexer.NewToken("(", lexer.OParen, 1, 1),
				lexer.NewToken("4", lexer.Number, 2, 1),
				lexer.NewToken("+", lexer.BinOperator, 3, 1),
				lexer.NewToken("2", lexer.Number, 4, 1),
				lexer.NewToken(")", lexer.CParen, 5, 1),
				lexer.NewToken("*", lexer.BinOperator, 6, 1),
				lexer.NewToken("3", lexer.Number, 7, 1),
				lexer.NewToken("<EOF>", lexer.EOF, 8, 0),
			},
			wantErr: false,
		},
		{
			name:  "Simple Lexical Tokenization",
			input: "let abc = 'hi'",
			want: []lexer.Token{
				lexer.NewToken("let", lexer.Let, 1, 3),
				lexer.NewToken("abc", lexer.Identifier, 5, 3),
				lexer.NewToken("=", lexer.Equals, 9, 1),
				lexer.NewToken("'hi'", lexer.String, 11, 4),
				lexer.NewToken("<EOF>", lexer.EOF, 15, 0),
			},
			wantErr: false,
		},

		{
			name:  "Comprehensive Operators and Numbers",
			input: "a = (5. * .2) + 1; // Test\n b != -1.5",
			want: []lexer.Token{
				lexer.NewToken("a", lexer.Identifier, 1, 1),
				lexer.NewToken("=", lexer.Equals, 3, 1),
				lexer.NewToken("(", lexer.OParen, 5, 1),
				lexer.NewToken("5", lexer.Number, 6, 1),
				lexer.NewToken(".", lexer.Dot, 7, 1),
				lexer.NewToken("*", lexer.BinOperator, 9, 1),
				lexer.NewToken(".2", lexer.Number, 11, 2),
				lexer.NewToken(")", lexer.CParen, 13, 1),
				lexer.NewToken("+", lexer.BinOperator, 15, 1),
				lexer.NewToken("1", lexer.Number, 17, 1),
				lexer.NewToken(";", lexer.SemiColon, 18, 1),
				lexer.NewToken("b", lexer.Identifier, 29, 1),
				lexer.NewToken("!=", lexer.ComOperator, 31, 2),
				lexer.NewToken("-", lexer.BinOperator, 34, 1),
				lexer.NewToken("1.5", lexer.Number, 35, 3),
				lexer.NewToken("<EOF>", lexer.EOF, 38, 0),
			},
			wantErr: false,
		},
		{
			name:  "Multi-line Comment and Whitespace",
			input: "  let x = 10; /* Multi\n line \n comment */ const y = 20;",
			want: []lexer.Token{
				lexer.NewToken("let", lexer.Let, 3, 3),
				lexer.NewToken("x", lexer.Identifier, 7, 1),
				lexer.NewToken("=", lexer.Equals, 9, 1),
				lexer.NewToken("10", lexer.Number, 11, 2),
				lexer.NewToken(";", lexer.SemiColon, 13, 1),
				lexer.NewToken("const", lexer.Const, 43, 5),
				lexer.NewToken("y", lexer.Identifier, 49, 1),
				lexer.NewToken("=", lexer.Equals, 51, 1),
				lexer.NewToken("20", lexer.Number, 53, 2),
				lexer.NewToken(";", lexer.SemiColon, 55, 1),
				lexer.NewToken("<EOF>", lexer.EOF, 56, 0),
			},
			wantErr: false,
		},
		{
			name:  "String Variations",
			input: `a = "double"; b = ''; c = "with space";`,
			want: []lexer.Token{
				lexer.NewToken("a", lexer.Identifier, 1, 1),
				lexer.NewToken("=", lexer.Equals, 3, 1),
				lexer.NewToken(`"double"`, lexer.String, 5, 8),
				lexer.NewToken(";", lexer.SemiColon, 13, 1),
				lexer.NewToken("b", lexer.Identifier, 15, 1),
				lexer.NewToken("=", lexer.Equals, 17, 1),
				lexer.NewToken("''", lexer.String, 19, 2),
				lexer.NewToken(";", lexer.SemiColon, 21, 1),
				lexer.NewToken("c", lexer.Identifier, 23, 1),
				lexer.NewToken("=", lexer.Equals, 25, 1),
				lexer.NewToken(`"with space"`, lexer.String, 27, 12),
				lexer.NewToken(";", lexer.SemiColon, 39, 1),
				lexer.NewToken("<EOF>", lexer.EOF, 40, 0),
			},
			wantErr: false,
		},
		{
			name:  "Identifier Variations",
			input: `let _a = $b1; try c$ = 1_000;`,
			want: []lexer.Token{
				lexer.NewToken("let", lexer.Let, 1, 3),
				lexer.NewToken("_a", lexer.Identifier, 5, 2),
				lexer.NewToken("=", lexer.Equals, 8, 1),
				lexer.NewToken("$b1", lexer.Identifier, 10, 3),
				lexer.NewToken(";", lexer.SemiColon, 13, 1),
				lexer.NewToken("try", lexer.Try, 15, 3),
				lexer.NewToken("c$", lexer.Identifier, 19, 2),
				lexer.NewToken("=", lexer.Equals, 22, 1),
				lexer.NewToken("1", lexer.Number, 24, 1),
				lexer.NewToken("_000", lexer.Identifier, 25, 4),
				lexer.NewToken(";", lexer.SemiColon, 29, 1),
				lexer.NewToken("<EOF>", lexer.EOF, 30, 0),
			},
			wantErr: false,
		},
		{
			name:  "Edge_Case_Numbers",
			input: "x = 0; y = .5; z = 5.;",
			want: []lexer.Token{
				lexer.NewToken("x", lexer.Identifier, 1, 1),
				lexer.NewToken("=", lexer.Equals, 3, 1),
				lexer.NewToken("0", lexer.Number, 5, 1),
				lexer.NewToken(";", lexer.SemiColon, 6, 1),
				lexer.NewToken("y", lexer.Identifier, 8, 1),
				lexer.NewToken("=", lexer.Equals, 10, 1),
				lexer.NewToken(".5", lexer.Number, 12, 2),
				lexer.NewToken(";", lexer.SemiColon, 14, 1),
				lexer.NewToken("z", lexer.Identifier, 16, 1),
				lexer.NewToken("=", lexer.Equals, 18, 1),
				lexer.NewToken("5", lexer.Number, 20, 1),
				lexer.NewToken(".", lexer.Dot, 21, 1),
				lexer.NewToken(";", lexer.SemiColon, 22, 1),
				lexer.NewToken("<EOF>", lexer.EOF, 23, 0),
			},
			wantErr: false,
		},
		{
			name:    "Empty Input",
			input:   "",
			want:    []lexer.Token{lexer.NewToken("<EOF>", lexer.EOF, 1, 0)},
			wantErr: false,
		},
		{
			name:    "Whitespace Only Input",
			input:   "  \n\t \r ",
			want:    []lexer.Token{lexer.NewToken("<EOF>", lexer.EOF, 8, 0)},
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
				lexer.NewToken("x", lexer.Identifier, 1, 1),
				lexer.NewToken("=", lexer.Equals, 3, 1),
				lexer.NewToken("0.123", lexer.Number, 5, 5),
				lexer.NewToken("<EOF>", lexer.EOF, 10, 0),
			},
			wantErr: false,
		},
		{
			name:  "Floating Point Number without Leading Zero",
			input: "y = .456",
			want: []lexer.Token{
				lexer.NewToken("y", lexer.Identifier, 1, 1),
				lexer.NewToken("=", lexer.Equals, 3, 1),
				lexer.NewToken(".456", lexer.Number, 5, 4),
				lexer.NewToken("<EOF>", lexer.EOF, 9, 0),
			},
			wantErr: false,
		},
		{
			name:  "Floating Point Number with Trailing Dot",
			input: "z = 789.",
			want: []lexer.Token{
				lexer.NewToken("z", lexer.Identifier, 1, 1),
				lexer.NewToken("=", lexer.Equals, 3, 1),
				lexer.NewToken("789", lexer.Number, 5, 3),
				lexer.NewToken(".", lexer.Dot, 8, 1),
				lexer.NewToken("<EOF>", lexer.EOF, 9, 0),
			},
			wantErr: false,
		},
		{
			name:  "Floating Point Number in Expression",
			input: "a = 1.2 + 3.4",
			want: []lexer.Token{
				lexer.NewToken("a", lexer.Identifier, 1, 1),
				lexer.NewToken("=", lexer.Equals, 3, 1),
				lexer.NewToken("1.2", lexer.Number, 5, 3),
				lexer.NewToken("+", lexer.BinOperator, 9, 1),
				lexer.NewToken("3.4", lexer.Number, 11, 3),
				lexer.NewToken("<EOF>", lexer.EOF, 14, 0),
			},
			wantErr: false,
		},
		{
			name:    "Floating Point Number with Multiple Dots Error",
			input:   "let b = 1.2.3",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := lexer.Tokenize(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Tokenize() mismatch:\nexpected: %v\n     got: %v", tt.want, got)
				}
			}
		})
	}
}
