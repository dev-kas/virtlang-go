package parser

import "github.com/dev-kas/virtlang-go/v4/lexer"

type Parser struct {
	tokens   []lexer.Token
	filename string
}

func New(filename string) *Parser {
	return &Parser{
		filename: filename,
		tokens:   []lexer.Token{},
	}
}
