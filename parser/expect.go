package parser

import (
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) expect(type_ lexer.TokenType) (*lexer.Token, *errors.SyntaxError) {
	prev := p.advance()
	if prev == nil || prev.Type != type_ {
		return nil, &errors.SyntaxError{
			Expected:   lexer.Stringify(type_),
			Got:        lexer.Stringify(prev.Type),
			Start:      prev.Start,
			Difference: prev.Difference,
		}
	}

	return prev, nil
}
