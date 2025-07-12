package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseUnaryExpr() (ast.Expr, *errors.SyntaxError) {
	tok := p.at()

	// Logical NOT
	if tok.Type == lexer.LogicalOperator && tok.Literal == "!" {
		p.advance()
		expr, err := p.parseUnaryExpr()
		if err != nil {
			return nil, err
		}

		return &ast.LogicalExpr{
			Operator: ast.LogicalNOT,
			LHS:      nil,
			RHS:      expr,
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   tok.StartLine,
				StartColumn: tok.StartCol,
				EndLine:     expr.GetSourceMetadata().EndLine,
				EndColumn:   expr.GetSourceMetadata().EndColumn,
			},
		}, nil
	}

	return p.parseCallMemberExpr()
}
