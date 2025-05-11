package parser

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/lexer"
)

func (p *Parser) parseStmt() (ast.Stmt, *errors.SyntaxError) {
	switch p.at().Type {
	case lexer.Let, lexer.Const:
		return p.parseVarDecl()
	case lexer.Fn:
		return p.parseFnDecl()
	case lexer.If:
		return p.parseIfStmt()
	case lexer.Class:
		return p.parseClass()
	default:
		return p.parseExpr()
	}
}
