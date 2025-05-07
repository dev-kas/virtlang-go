package parser

import (
	"github.com/dev-kas/virtlang-go/ast"
	"github.com/dev-kas/virtlang-go/errors"
)

func (p *Parser) parseBreakStmt() (ast.Expr, *errors.SyntaxError) {
	p.advance() // break

	return &ast.BreakStmt{}, nil
}
