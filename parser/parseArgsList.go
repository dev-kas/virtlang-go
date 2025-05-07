package parser

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/lexer"
)

func (p *Parser) parseArgsList() ([]ast.Expr, *errors.SyntaxError) {
	var args = make([]ast.Expr, 1)

	arg, err := p.parseAssignmentExpr()
	if err != nil {
		return nil, err
	}
	args[0] = arg

	for p.at().Type == lexer.Comma {
		p.advance()

		arg, err := p.parseAssignmentExpr()
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return args, nil
}
