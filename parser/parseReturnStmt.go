package parser

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/errors"
)

func (p *Parser) parseReturnStmt() (ast.Expr, *errors.SyntaxError) {
	p.advance() // return
	var value ast.Expr

	if p.isEOF() {
		value = nil
	} else {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		value = expr
	}

	return &ast.ReturnStmt{
		Value: value,
	}, nil
}
