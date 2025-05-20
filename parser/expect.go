package parser

import (
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/lexer"
)

func (p *Parser) expect(type_ lexer.TokenType) (*lexer.Token, *errors.SyntaxError) {
	prev := p.advance()
	if prev == nil || prev.Type != type_ {
		return nil, &errors.SyntaxError{
			Expected: lexer.Stringify(type_),
			Got:      lexer.Stringify(prev.Type),
			Start:    errors.Position{Line: prev.StartLine, Col: prev.StartCol},
			End:      errors.Position{Line: prev.EndLine, Col: prev.EndCol},
		}
	}

	return prev, nil
}
