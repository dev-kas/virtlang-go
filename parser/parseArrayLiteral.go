package parser

import (
	"VirtLang/ast"
	"VirtLang/errors"
	"VirtLang/lexer"
	"strconv"
)

func (p *Parser) parseArrayLiteral() (ast.Expr, *errors.SyntaxError) {
	p.advance() // [

	elements := []ast.Expr{}

	for !p.isEOF() && p.at().Type != lexer.CBracket {
		element, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)

		if p.at().Type == lexer.Comma {
			p.advance()
		}
	}

	_, err := p.expect(lexer.CBracket)
	if err != nil {
		return nil, err
	}

	properties := []ast.Property{}

	for index, element := range elements {
		properties = append(properties, ast.Property{
			Key:   strconv.Itoa(index),
			Value: element,
		})
	}

	return &ast.ObjectLiteral{
		Properties: properties,
	}, nil
}
