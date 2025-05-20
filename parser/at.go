package parser

import "github.com/dev-kas/virtlang-go/v3/lexer"

func (p *Parser) at() lexer.Token {
	return p.tokens[0]
}
