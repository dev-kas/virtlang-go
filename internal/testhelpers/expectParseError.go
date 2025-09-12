package testhelpers

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v4/parser"
)

// ExpectParseError checks that parsing fails for invalid sources
func ExpectParseError(t *testing.T, src string) {
	t.Helper()
	p := parser.New("test")
	_, err := p.ProduceAST(src)
	if err == nil {
		t.Fatal("Expected parse error")
	}
}
