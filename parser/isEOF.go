package parser

import "github.com/dev-kas/virtlang-go/lexer"

func (p *Parser) isEOF() bool {
	return p.at().Type == lexer.EOF
}
