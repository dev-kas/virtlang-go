package parser

import (
	"VirtLang/ast"
	"VirtLang/errors"
	// "VirtLang/lexer"
)

func (p *Parser) parseExpr() (ast.Expr, *errors.SyntaxError) {
	// if (p.at().Type == lexer.Fn) { // For Immediately Invoked Function Expression
	// 	return p.parseFnDecl()
	// }

	return p.parseAssignmentExpr()
}
