package parser

import (
	"VirtLang/ast"
	"VirtLang/lexer"
	"VirtLang/errors"
)

func (p *Parser) parseWhileLoop() (ast.Expr, *errors.SyntaxError) {
	p.advance() // while

	p.expect(lexer.OParen)
	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	p.expect(lexer.CParen)

	p.expect(lexer.OBrace)
	body := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	p.expect(lexer.CBrace)

	return &ast.WhileLoop{
		Condition: condition,
		Body:      body,
	}, nil
}