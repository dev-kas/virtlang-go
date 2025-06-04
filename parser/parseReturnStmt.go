package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
)

func (p *Parser) parseReturnStmt() (ast.Expr, *errors.SyntaxError) {
	start := p.advance() // return
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
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}, nil
}
