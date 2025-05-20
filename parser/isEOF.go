package parser

import "github.com/dev-kas/virtlang-go/v3/lexer"

func (p *Parser) isEOF() bool {
	return p.at().Type == lexer.EOF
}
