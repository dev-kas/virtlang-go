package parser

import "VirtLang/lexer"

func (p *Parser) at() lexer.Token {
	return p.tokens[0]
}
