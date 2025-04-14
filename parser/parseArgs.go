package parser

import (
	"VirtLang/ast"
	"VirtLang/lexer"
	"VirtLang/errors"
)

func (p *Parser) parseArgs() ([]ast.Expr, *errors.SyntaxError) {
	p.expect(lexer.OParen)
	var args []ast.Expr
	if p.at().Type == lexer.CParen {
		args = []ast.Expr{}
	} else {
		newArgs, err := p.parseArgsList()
		if err != nil {
			return nil, err
		}
		args = newArgs
	}
	p.expect(lexer.CParen)
	return args, nil
}
