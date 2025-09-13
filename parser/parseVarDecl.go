package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseVarDecl() (ast.Declaration, *errors.SyntaxError) {
	start := p.at()
	isConstant := p.advance().Type == lexer.Const

	if p.at().Type != lexer.Identifier {
		// can be either a destructure pattern or an invalid token
		pattern, err := p.parseDestructurePattern()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(lexer.Equals)
		if err != nil {
			return nil, err
		}

		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		return &ast.DestructureDeclaration{
			Pattern:  pattern,
			Constant: isConstant,
			Value:    value,
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}, nil
	}

	ident, err := p.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	name := ident.Literal
	_, err = p.expect(lexer.Equals)
	if err != nil {
		return nil, err
	}

	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.VarDeclaration{
		Identifier: name,
		Value:      value,
		Constant:   isConstant,
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}, nil
}
