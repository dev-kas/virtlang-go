package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
)

func (p *Parser) parseContinueStmt() (ast.Expr, *errors.SyntaxError) {
	p.advance() // break

	return &ast.ContinueStmt{}, nil
}
