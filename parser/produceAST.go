package parser

import (
	"github.com/dev-kas/virtlang-go/ast"
	// "github.com/dev-kas/virtlang-go/errors"
	"github.com/dev-kas/virtlang-go/lexer"
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
