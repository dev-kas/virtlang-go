package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseMemberExpr() (ast.Expr, *errors.SyntaxError) {
	start := p.at()
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
					Expected: "Identifier",
					Got:      property.GetType().String(),
					Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
					End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
				}
			}
		} else {
			computed = true
			property, err = p.parseExpr()
			if err != nil {
				return nil, err
			}

			_, err = p.expect(lexer.CBracket)
			if err != nil {
				return nil, err
			}
		}

		obj = &ast.MemberExpr{
			Object:   obj,
			Value:    property,
			Computed: computed,
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}
	}

	return obj, nil
}
