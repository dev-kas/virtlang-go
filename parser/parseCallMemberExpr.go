package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
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
