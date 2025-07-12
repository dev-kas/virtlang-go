package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/lexer"
)

func (p *Parser) parseLogicalExpr() (ast.Expr, *errors.SyntaxError) {
	start := p.at()
	lhs, err := p.parseComparisonExpr()
	if err != nil {
		return nil, err
	}

	if p.at().Type == lexer.LogicalOperator {
		operatorToken := p.advance().Literal

		var operator ast.LogicalOperator
		switch operatorToken {
		case "&&":
			operator = ast.LogicalAND
		case "||":
			operator = ast.LogicalOR
		case "??":
			operator = ast.LogicalNilCoalescing
		case "!":
			operator = ast.LogicalNOT
		default:
			return nil, &errors.SyntaxError{
				Expected: "&&, ||, ??, !",
				Got:      operatorToken,
				Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
				End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
			}
		}

		rhs, err := p.parseLogicalExpr()
		if err != nil {
			return nil, err
		}

		return &ast.LogicalExpr{
			LHS:      &lhs,
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
