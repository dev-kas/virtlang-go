package parser

import (
	"VirtLang/ast"
	"VirtLang/errors"
	"VirtLang/lexer"
)

func (p *Parser) parseIfStmt() (*ast.IfStatement, *errors.SyntaxError) {
	p.advance() // if

	if _, err := p.expect(lexer.OParen); err != nil {
		return nil, err
	}

	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.CParen); err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.OBrace); err != nil {
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

	if _, err := p.expect(lexer.CBrace); err != nil {
		return nil, err
	}

	return &ast.IfStatement{
		Condition: condition,
		Body:      body,
	}, nil
}
