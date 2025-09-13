package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseDestructureObjectPattern() (ast.DestructurePattern, *errors.SyntaxError) {
	start, err := p.expect(lexer.OBrace) // {
	if err != nil {
		return nil, err
	}

	props := []ast.DestructureObjectProperty{}
	var rest *string

	for p.at().Type != lexer.CBrace && !p.isEOF() {
		isRest := false
		// could be ...
		if p.at().Type == lexer.Dot {
			for i := 0; i < 3; i++ {
				_, err := p.expect(lexer.Dot)
				if err != nil {
					return nil, err
				}
			}

			isRest = true
		}
		// identifier
		ident, err := p.expect(lexer.Identifier)
		if err != nil {
			return nil, err
		}

		if rest != nil {
			return nil, &errors.SyntaxError{
				Expected: "rest element to be last element",
				Got:      ident.Literal,
				Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
				End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
			}
		}

		prop := &ast.DestructureObjectProperty{
			Key:  ident.Literal,
			Name: ident.Literal,
		}

		// either `:` for renaming or nested destruct
		// or `=` for default value
		// or `,` for further destructuring

		if p.at().Type == lexer.Colon {
			// renaming or nested destructuring
			_, err := p.expect(lexer.Colon)
			if err != nil {
				return nil, err
			}

			if p.at().Type == lexer.Identifier {
				// renaming
				ident, err := p.expect(lexer.Identifier)
				if err != nil {
					return nil, err
				}

				prop.Name = ident.Literal
			} else {
				// nested destructure
				children, err := p.parseDestructurePattern()
				if err != nil {
					return nil, err
				}

				prop.DeconstructChildren = children
			}
		}

		if p.at().Type == lexer.Equals {
			// default value
			_, err := p.expect(lexer.Equals)
			if err != nil {
				return nil, err
			}

			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}

			prop.Default = expr
		}

		if isRest {
			rest = &ident.Literal
		} else {
			props = append(props, *prop)
		}


		if p.at().Type == lexer.Comma {
			// further destructuring
			_, err := p.expect(lexer.Comma)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	_, err = p.expect(lexer.CBrace) // }
	if err != nil {
		return nil, err
	}

	return &ast.DestructureObjectPattern{
		Properties: props,
		Rest:       rest,
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}, nil
}
