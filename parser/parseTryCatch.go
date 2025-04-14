package parser

import (
	"VirtLang/ast"
	"VirtLang/lexer"
	"VirtLang/errors"
)

func (p *Parser) parseTryCatch() (ast.Expr, *errors.SyntaxError) {
	p.advance() // try
	p.expect(lexer.OBrace)

	tryBody := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		tryBody = append(tryBody, stmt)
	}

	p.expect(lexer.CBrace)
	p.expect(lexer.Catch)

	cVar, err := p.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	p.expect(lexer.OBrace)

	catchBody := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		catchBody = append(catchBody, stmt)
	}

	p.expect(lexer.CBrace)

	return &ast.TryCatchStmt{
		Try:  tryBody,
		CatchVar: cVar.Literal,
		Catch: catchBody,
	}, nil
}