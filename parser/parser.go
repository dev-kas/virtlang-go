package parser

import "github.com/dev-kas/virtlang-go/v2/lexer"

// type Parser interface {
// 	at() lexer.Token
// 	advance() lexer.Token
// 	expect() (lexer.Token, error)
// 	isEOF() bool
// 	ProduceAST() (ast.Program, error)
// 	parseStmt() (ast.Stmt, error)
// 	parseIfStmt() (ast.Stmt, error)
// 	parseFnDecl() (ast.Stmt, error)
// 	parseVarDecl() (ast.Stmt, error)
// 	parseTryCatch() (ast.Expr, error)
// 	parseExpr() (ast.Expr, error)
// 	parseAssignmentExpr() (ast.Expr, error)
// 	parseComparisonExpr() (ast.Expr, error)
// 	parseObjectExpr() (ast.Expr, error)
// 	parseAdditiveExpr() (ast.Expr, error)
// 	parseMultiplicativeExpr() (ast.Expr, error)
// 	parseCallMemberExpr() (ast.Expr, error)
// 	parseCallExpr() (ast.Expr, error)
// 	parseArgs() ([]ast.Expr, error)
// 	parseArgsList() ([]ast.Expr, error)
// 	parseMemberExpr() (ast.Expr, error)
// 	parseArrayLiteral() (ast.Expr, error)
// 	parseWhileLoop() (ast.Expr, error)
// 	parseReturnStmt() (ast.Expr, error)
// 	parsePrimaryExpr() (ast.Expr, error)
// }

type Parser struct {
	tokens []lexer.Token
}

func New() *Parser {
	return &Parser{
		tokens: []lexer.Token{},
	}
}
