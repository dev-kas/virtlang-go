package parser

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	"github.com/dev-kas/virtlang-go/v2/errors"
)

func (p *Parser) parseMultiplicativeExpr() (ast.Expr, *errors.SyntaxError) {
	lhs, err := p.parseCallMemberExpr()
	if err != nil {
		return nil, err
	}

	for p.at().Literal == "*" || p.at().Literal == "/" || p.at().Literal == "%" {
		operatorLiteral := p.advance().Literal
		var operator ast.BinaryOperator
		switch operatorLiteral {
		case "*":
			operator = ast.Multiply
		case "/":
			operator = ast.Divide
		case "%":
			operator = ast.Modulo
		default:
			return nil, &errors.SyntaxError{
				Expected:   "+, -",
				Got:        operatorLiteral,
				Start:      p.at().Start,
				Difference: p.at().Difference,
			}
		}

		rhs, err := p.parseCallMemberExpr()
		if err != nil {
			return nil, err
		}
		lhs = &ast.BinaryExpr{
			Operator: operator,
			LHS:      lhs,
			RHS:      rhs,
		}
	}

	return lhs, nil
}
