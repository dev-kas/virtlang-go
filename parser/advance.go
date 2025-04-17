package parser

import "github.com/dev-kas/VirtLang-Go/lexer"

func (p *Parser) advance() *lexer.Token {
	if len(p.tokens) == 0 {
		return nil
	}
	prev := p.tokens[0]
	p.tokens = p.tokens[1:]
	return &prev
}
