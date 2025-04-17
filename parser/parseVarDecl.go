package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) parseVarDecl() (*ast.VarDeclaration, *errors.SyntaxError) {
	isConstant := p.advance().Type == lexer.Const
	ident, err := p.expect(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	name := ident.Literal
	_, err = p.expect(lexer.Equals)
	if err != nil {
		return nil, err
	}

	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.VarDeclaration{
		Identifier: name,
		Value:      value,
		Constant:   isConstant,
	}, nil
}
