package parser

import "github.com/dev-kas/virtlang-go/v2/lexer"

func (p *Parser) at() lexer.Token {
	return p.tokens[0]
}
