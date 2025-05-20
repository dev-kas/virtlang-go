package parser

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/lexer"
)

func (p *Parser) parseComparisonExpr() (ast.Expr, *errors.SyntaxError) {
	start := p.at()
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
				Expected: "<, ==, >, !=, <=, >=",
				Got:      operatorToken,
				Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
				End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
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
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}, nil
	}

	return lhs, nil
}
