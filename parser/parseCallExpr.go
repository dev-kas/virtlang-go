package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseCallExpr(callee ast.Expr) (*ast.CallExpr, *errors.SyntaxError) {
	start := p.at()
	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}

	callExpr := &ast.CallExpr{
		Callee: callee,
		Args:   args,
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}

	if p.at().Type == lexer.OParen {
		expr, err := p.parseCallExpr(callExpr)
		if err != nil {
			return nil, err
		}
		callExpr = expr
	}

	return callExpr, nil
}
