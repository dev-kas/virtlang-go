package parser

import (
	"VirtLang/ast"
	"VirtLang/errors"
	"VirtLang/lexer"
)

func (p *Parser) parseObjectExpr() (ast.Expr, *errors.SyntaxError) {
	if p.at().Type != lexer.OBrace {
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
	}, nil
}
