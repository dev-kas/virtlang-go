package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
)

func (p *Parser) parseExpr() (ast.Expr, *errors.SyntaxError) {
	if p.at().Type == lexer.Fn { // For Immediately Invoked Function Expression
		return p.parseFnDecl()
	}

	return p.parseAssignmentExpr()
}
