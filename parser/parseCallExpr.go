package parser

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/lexer"
)

func (p *Parser) parseCallExpr(callee ast.Expr) (*ast.CallExpr, *errors.SyntaxError) {
	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}

	callExpr := &ast.CallExpr{
		Callee: callee,
		Args:   args,
	}

	if p.at().Type == lexer.OParen {
		expr, err := p.parseCallExpr(callExpr)
		if err != nil {
			return nil, err
		}
		callExpr = expr
	}

	return callExpr, nil
}
