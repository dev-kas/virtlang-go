package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) parseAssignmentExpr() (ast.Expr, *errors.SyntaxError) {
	lhs, err := p.parseComparisonExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == lexer.Equals {
		p.advance()
		value, err := p.parseAssignmentExpr()
		if err != nil {
			return nil, err
		}
		return &ast.VarAssignmentExpr{
			Value:    value,
			Assignee: lhs,
		}, nil
	}

	return lhs, nil
}
