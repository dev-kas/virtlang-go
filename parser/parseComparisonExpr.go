package parser

import (
	"VirtLang/ast"
	"VirtLang/errors"
	"VirtLang/lexer"
)

func (p *Parser) parseComparisonExpr() (ast.Expr, *errors.SyntaxError) {
	lhs, err := p.parseObjectExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == lexer.ComOperator {
		operatorToken := p.advance().Literal

		var operator ast.CompareOperator
		switch operatorToken {
		case "<":
			operator = ast.LessThan
		case "==":
			operator = ast.Equal
		case ">":
			operator = ast.GreaterThan
		case "!=":
			operator = ast.NotEqual
		case "<=":
			operator = ast.LessThanEqual
		case ">=":
			operator = ast.GreaterThanEqual
		default:
			return nil, &errors.SyntaxError{
				Expected:   "<, ==, >, !=, <=, >=",
				Got:        operatorToken,
				Start:      p.at().Start,
				Difference: p.at().Difference,
			}
		}

		rhs, err := p.parseObjectExpr()
		if err != nil {
			return nil, err
		}

		return &ast.CompareExpr{
			LHS:      lhs,
			RHS:      rhs,
			Operator: operator,
		}, nil
	}

	return lhs, nil
}
