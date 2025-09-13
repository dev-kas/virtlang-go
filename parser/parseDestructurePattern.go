package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseDestructurePattern() (ast.DestructurePattern, *errors.SyntaxError) {
	// destructure pattern starts with either { (for object based patterns) or [ (for array based patterns)

	if p.at().Type == lexer.OBracket {
		return p.parseDestructureArrayPattern()
	}

	return p.parseDestructureObjectPattern()
}
