package parser

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/lexer"
)

func (p *Parser) parseExpr() (ast.Expr, *errors.SyntaxError) {
	if p.at().Type == lexer.Fn { // For Immediately Invoked Function Expression
		return p.parseFnDecl()
	}

	return p.parseAssignmentExpr()
}
