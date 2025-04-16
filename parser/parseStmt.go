package parser

import (
	"VirtLang/ast"
	"VirtLang/errors"
	"VirtLang/lexer"
)

func (p *Parser) parseStmt() (ast.Stmt, *errors.SyntaxError) {
	switch p.at().Type {
	case lexer.Let, lexer.Const:
		return p.parseVarDecl()
	case lexer.Fn:
		return p.parseFnDecl()
	case lexer.If:
		return p.parseIfStmt()
	default:
		return p.parseExpr()
	}
}
