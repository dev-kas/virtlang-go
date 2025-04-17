package parser

import "github.com/dev-kas/VirtLang-Go/lexer"

func (p *Parser) at() lexer.Token {
	return p.tokens[0]
}
