package parser

import (
	"github.com/dev-kas/virtlang-go/ast"
	"github.com/dev-kas/virtlang-go/errors"
	"github.com/dev-kas/virtlang-go/lexer"
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
