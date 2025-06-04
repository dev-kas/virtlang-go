package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseTryCatch() (ast.Expr, *errors.SyntaxError) {
	start := p.advance() // try
	if _, err := p.expect(lexer.OBrace); err != nil {
		return nil, err
	}

	tryBody := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		tryBody = append(tryBody, stmt)
	}

	if _, err := p.expect(lexer.CBrace); err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.Catch); err != nil {
		return nil, err
	}

	cVar, err := p.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.OBrace); err != nil {
		return nil, err
	}

	catchBody := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		catchBody = append(catchBody, stmt)
	}

	if _, err := p.expect(lexer.CBrace); err != nil {
		return nil, err
	}

	return &ast.TryCatchStmt{
		Try:      tryBody,
		CatchVar: cVar.Literal,
		Catch:    catchBody,
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}, nil
}
