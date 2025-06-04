package parser

import (
	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/errors"
)

func (p *Parser) parseAdditiveExpr() (ast.Expr, *errors.SyntaxError) {
	start := p.at()
	lhs, err := p.parseMultiplicativeExpr()
	if err != nil {
		return nil, err
	}

	for p.at().Literal == "+" || p.at().Literal == "-" {
		operatorLiteral := p.advance().Literal
		var operator ast.BinaryOperator
		switch operatorLiteral {
		case "+":
			operator = ast.Plus
		case "-":
			operator = ast.Minus
		default:
			return nil, &errors.SyntaxError{
				Expected: "+, -",
				Got:      operatorLiteral,
				Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
				End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
			}
		}
		rhs, err := p.parseMultiplicativeExpr()
		if err != nil {
			return nil, err
		}
		lhs = &ast.BinaryExpr{
			Operator: operator,
			LHS:      lhs,
			RHS:      rhs,
			SourceMetadata: ast.SourceMetadata{
				Filename:    p.filename,
				StartLine:   start.StartLine,
				StartColumn: start.StartCol,
				EndLine:     p.at().EndLine,
				EndColumn:   p.at().EndCol,
			},
		}
	}

	return lhs, nil
}
