package parser

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/lexer"
)

func (p *Parser) parseClassStmt() (ast.Stmt, *errors.SyntaxError) {
	start := p.at()
	isPrivate := p.at().Type == lexer.Private
	if isPrivate {
		p.advance()
	} else {
		_, err := p.expect(lexer.Public)
		if err != nil {
			return nil, err
		}
	}

	ident, err := p.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	name := ident.Literal

	isFunc := p.at().Type == lexer.OParen
	if isFunc {
		args, err := p.parseArgs()
		if err != nil {
			return nil, err
		}
		params := []string{}

		for _, arg := range args {
			if arg.GetType() != ast.IdentifierNode {
				return nil, &errors.SyntaxError{
					Expected: "Identifier",
					Got:      arg.GetType().String(),
					Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
					End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
				}
			}
			params = append(params, arg.(*ast.Identifier).Symbol)
		}

		_, err = p.expect(lexer.OBrace)
		if err != nil {
			return nil, err
		}
		body := []ast.Stmt{}
		for !p.isEOF() && p.at().Type != lexer.CBrace {
			stmt, err := p.parseStmt()
			if err != nil {
				return nil, err
			}
			body = append(body, stmt)
		}
		_, err = p.expect(lexer.CBrace)
		if err != nil {
			return nil, err
		}
		return &ast.ClassMethod{
			Name:     name,
			Body:     body,
			Params:   params,
			IsPublic: !isPrivate,
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}, nil
	} else {
		at := p.at()
		if at.Type != lexer.Equals {
			return &ast.ClassProperty{
				Name:     name,
				Value:    nil,
				IsPublic: !isPrivate,
			}, nil
		}
		p.advance()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.ClassProperty{
			Name:     name,
			Value:    value,
			IsPublic: !isPrivate,
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}, nil
	}
}
