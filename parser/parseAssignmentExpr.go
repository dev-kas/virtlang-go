package parser

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/lexer"
)

func (p *Parser) parseAssignmentExpr() (ast.Expr, *errors.SyntaxError) {
	start := p.at()
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
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}, nil
	}

	return lhs, nil
}
