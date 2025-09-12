package testhelpers

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/parser"
)

// MustParse parses source code and returns AST program or fails the test
func MustParse(t *testing.T, src string) *ast.Program {
	t.Helper()
	p := parser.New("test")
	prog, err := p.ProduceAST(src)
	if err != nil {
		t.Fatal(err)
	}
	return prog
}
