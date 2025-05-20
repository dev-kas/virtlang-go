package parser

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/lexer"
)

func (p *Parser) parseArrayLiteral() (ast.Expr, *errors.SyntaxError) {
	start := p.advance() // [

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

	// properties := []ast.Property{}

	// for index, element := range elements {
	// 	properties = append(properties, ast.Property{
	// 		Key:   strconv.Itoa(index),
	// 		Value: element,
	// 	})
	// }

	return &ast.ArrayLiteral{
		Elements: elements,
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}, nil
}
