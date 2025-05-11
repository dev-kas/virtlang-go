package parser

import "github.com/dev-kas/virtlang-go/v2/lexer"

type Parser struct {
	tokens []lexer.Token
}

func New() *Parser {
	return &Parser{
		tokens: []lexer.Token{},
	}
}
