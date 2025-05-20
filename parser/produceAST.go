package parser

import (
	"github.com/dev-kas/virtlang-go/v3/ast"
	// "github.com/dev-kas/virtlang-go/v3/errors"
	"github.com/dev-kas/virtlang-go/v3/lexer"
)

func (p *Parser) ProduceAST(srcCode string) (*ast.Program, error) {
	tokens, err := lexer.Tokenize(srcCode)
	if err != nil {
		return nil, err
	}

	p.tokens = tokens

	program := ast.Program{
		Stmts: []ast.Stmt{},
		SourceMetadata: ast.SourceMetadata{
			Filename: srcCode,
			StartLine: 1,
			StartColumn: 1,
			EndLine: tokens[len(tokens)-1].EndLine,
			EndColumn: tokens[len(tokens)-1].EndCol,
		},
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
