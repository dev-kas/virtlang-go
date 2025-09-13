package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseDestructureArrayPattern() (ast.DestructurePattern, *errors.SyntaxError) {
	start, err := p.expect(lexer.OBracket) // [
	if err != nil {
		return nil, err
	}

	elems := []ast.DestructureArrayElement{}
	var rest *string

	for p.at().Type != lexer.CBrace && !p.isEOF() {
		// for element skipping. eg: [abc, , xyz]
		isSkip := p.at().Type == lexer.Comma

		isRest := false
		// could be ...
		if p.at().Type == lexer.Dot && !isSkip { // rest only possible if not skipped
			for i := 0; i < 3; i++ {
				_, err := p.expect(lexer.Dot)
				if err != nil {
					return nil, err
				}
			}

			isRest = true
		}

		var elem = ast.DestructureArrayElement{
			Skipped: isSkip,
		}

		// fill only required fields if skipped
		if isSkip {
			elem.SourceMetadata = ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			}
			elems = append(elems, elem)
			p.advance() // ,
			continue
		}

		// not skipped
		// first try to read name, else parse children destructure pattern
		if p.at().Type == lexer.Identifier {
			ident, err := p.expect(lexer.Identifier)
			if err != nil {
				return nil, err
			}

			elem.Name = ident.Literal
		} else {
			// it's a children destructure pattern
			pattern, err := p.parseDestructurePattern()
			if err != nil {
				return nil, err
			}
			elem.DeconstructChildren = pattern
		}

		// for elements after rest element
		if rest != nil {
			got := elem.Name
			if elem.DeconstructChildren != nil { // it is a children destructure pattern
				got = elem.DeconstructChildren.GetType().String()
			}
			return nil, &errors.SyntaxError{
				Expected: "rest element to be last element",
				Got:      got,
				Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
				End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
			}
		}

		// at this point, the only valid possibilities are:
		// either = for defaluting
		// or , for next element

		if p.at().Type == lexer.Equals {
			// for default value
			_, err := p.expect(lexer.Equals)
			if err != nil {
				return nil, err
			}

			val, err := p.parseExpr()
			if err != nil {
				return nil, err
			}

			elem.Default = val
		}

		// assign rest
		if isRest {
			rest = &elem.Name
		} else {
			elems = append(elems, elem)
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

	_, err = p.expect(lexer.CBracket) // ]
	if err != nil {
		return nil, err
	}

	return &ast.DestructureArrayPattern{
		Elements: elems,
		Rest:     rest,
		SourceMetadata: ast.SourceMetadata{
			Filename:    p.filename,
			StartLine:   start.StartLine,
			StartColumn: start.StartCol,
			EndLine:     p.at().EndLine,
			EndColumn:   p.at().EndCol,
		},
	}, nil
}
