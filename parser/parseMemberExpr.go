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

			var key *lexer.Token
			var err *errors.SyntaxError
			if _, ok := lexer.REVERSE_KEYWORDS[p.at().Type]; ok {
				key = p.advance()
			} else {
				key, err = p.expect(lexer.Identifier)
				if err != nil {
					return nil, err
				}
			}

			property = &ast.Identifier{
				Symbol: key.Literal,
				SourceMetadata: ast.SourceMetadata{
					Filename:    p.filename,
					StartLine:   start.StartLine,
					StartColumn: start.StartCol,
					EndLine:     p.at().EndLine,
					EndColumn:   p.at().EndCol,
				},
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
