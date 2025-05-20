package parser

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/lexer"
)

func (p *Parser) parseClass() (*ast.Class, *errors.SyntaxError) {
	start := p.advance() // class

	ident, err := p.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	name := ident.Literal
	// TODO: Implement class inheritance
	_, err = p.expect(lexer.OBrace)
	if err != nil {
		return nil, err
	}

	body := []ast.Stmt{}
	constructor := &ast.ClassMethod{
		Name: "constructor",
		Body: []ast.Stmt{},
	}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseClassStmt()
		if err != nil {
			return nil, err
		}
		if stmt.GetType() == ast.ClassMethodNode && stmt.(*ast.ClassMethod).Name == "constructor" {
			constructor = stmt.(*ast.ClassMethod)
		} else {
			body = append(body, stmt)
		}
	}

	_, err = p.expect(lexer.CBrace)
	if err != nil {
		return nil, err
	}

	return &ast.Class{
		Name:        name,
		Body:        body,
		Constructor: constructor,
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}, nil
}
