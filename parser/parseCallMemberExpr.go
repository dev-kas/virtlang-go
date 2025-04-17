package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) parseCallMemberExpr() (ast.Expr, *errors.SyntaxError) {
	member, err := p.parseMemberExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == lexer.OParen {
		parsedCallExpr, err := p.parseCallExpr(member)
		if err != nil {
			return nil, err
		}
		return parsedCallExpr, nil
	}

	return member, nil
}
