package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseObjectExpr() (ast.Expr, *errors.SyntaxError) {
	start := p.at()
	if start.Type != lexer.OBrace {
		return p.parseAdditiveExpr()
	}

	p.advance()

	properties := []ast.Property{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		key, err := p.expect(lexer.Identifier)
		if err != nil {
			return nil, err
		}

		if p.at().Type == lexer.Comma { // { key, }
			p.advance()
			properties = append(properties, ast.Property{
				Key:   key.Literal,
				Value: nil,
			})
			continue
		} else if p.at().Type == lexer.CBrace { // { key }
			properties = append(properties, ast.Property{
				Key:   key.Literal,
				Value: nil,
			})
			break
		}

		// { key: value }
		_, err = p.expect(lexer.Colon)
		if err != nil {
			return nil, err
		}

		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		properties = append(properties, ast.Property{
			Key:   key.Literal,
			Value: value,
		})

		if p.at().Type != lexer.CBrace {
			_, err = p.expect(lexer.Comma)
			if err != nil {
				return nil, err
			}
		}
	}

	_, err := p.expect(lexer.CBrace)
	if err != nil {
		return nil, err
	}

	return &ast.ObjectLiteral{
		Properties: properties,
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}, nil
}
