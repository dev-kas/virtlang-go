package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) parseMemberExpr() (ast.Expr, *errors.SyntaxError) {
	obj, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	for p.at().Type == lexer.Dot || p.at().Type == lexer.OBracket {
		operator := p.advance()
		var property ast.Expr
		var computed bool

		if operator.Type == lexer.Dot {
			computed = false
			property, err = p.parsePrimaryExpr()
			if err != nil {
				return nil, err
			}

			if property.GetType() != ast.IdentifierNode {
				return nil, &errors.SyntaxError{
					Expected:   "Identifier",
					Got:        property.GetType().String(),
					Start:      p.at().Start,
					Difference: p.at().Difference,
				}
			}
		} else {
			computed = true
			property, err = p.parseExpr()
			if err != nil {
				return nil, err
			}

			p.expect(lexer.CBracket)
		}

		obj = &ast.MemberExpr{
			Object:   obj,
			Value:    property,
			Computed: computed,
		}
	}

	return obj, nil
}
