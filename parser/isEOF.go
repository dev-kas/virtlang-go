package parser

import "github.com/dev-kas/VirtLang-Go/lexer"

func (p *Parser) isEOF() bool {
	return p.at().Type == lexer.EOF
}
