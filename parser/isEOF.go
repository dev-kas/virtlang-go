package parser

import "VirtLang/lexer"

func (p *Parser) isEOF() bool {
	return p.at().Type == lexer.EOF
}
