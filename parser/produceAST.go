package parser

import (
	"github.com/dev-kas/virtlang-go/v2/ast"
	// "github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/lexer"
)

func (p *Parser) ProduceAST(srcCode string) (*ast.Program, error) {
	tokens, err := lexer.Tokenize(srcCode)
	if err != nil {
		return nil, err
	}

	p.tokens = tokens

	program := ast.Program{
		Stmts: []ast.Stmt{},
	}

	for !p.isEOF() {
		parsed, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		program.Stmts = append(program.Stmts, parsed)
	}

	return &program, nil
}
