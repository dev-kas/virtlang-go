package parser

import "github.com/dev-kas/virtlang-go/lexer"

func (p *Parser) at() lexer.Token {
	return p.tokens[0]
}
