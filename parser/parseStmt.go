package parser

import (
	"github.com/dev-kas/VirtLang-Go/ast"
	"github.com/dev-kas/VirtLang-Go/errors"
	"github.com/dev-kas/VirtLang-Go/lexer"
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
