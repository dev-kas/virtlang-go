package parser

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/errors"
)

func (p *Parser) parseBreakStmt() (ast.Expr, *errors.SyntaxError) {
	p.advance() // break

	return &ast.BreakStmt{}, nil
}
