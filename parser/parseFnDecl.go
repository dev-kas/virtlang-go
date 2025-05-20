package parser

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	"github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/lexer"
)

func (p *Parser) parseFnDecl() (*ast.FnDeclaration, *errors.SyntaxError) {
	p.advance() // fn

	var name *lexer.Token

	if at := p.at(); at.Type == lexer.Identifier {
		name = p.advance()
	} else {
		token := lexer.NewToken("", lexer.Identifier, 0, 0, 0, 0)
		name = &token
	}
	isAnonymous := name.Type == lexer.OParen
	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}
	params := []string{}

	for _, arg := range args {
		if arg.GetType() != ast.IdentifierNode {
			return nil, &errors.SyntaxError{
				Expected: "Identifier",
				Got:      arg.GetType().String(),
				Start:    errors.Position{Line: p.at().StartLine, Col: p.at().StartCol},
				End:      errors.Position{Line: p.at().EndLine, Col: p.at().EndCol},
			}
		}
		params = append(params, arg.(*ast.Identifier).Symbol)
	}

	if _, err := p.expect(lexer.OBrace); err != nil {
		return nil, err
	}

	body := []ast.Stmt{}

	for !p.isEOF() && p.at().Type != lexer.CBrace {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	if _, err := p.expect(lexer.CBrace); err != nil {
		return nil, err
	}

	fname := ""

	if !isAnonymous {
		fname = name.Literal
	}

	return &ast.FnDeclaration{
		Params:    params,
		Name:      fname,
		Body:      body,
		Anonymous: isAnonymous,
	}, nil
}
