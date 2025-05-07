<div align="center">
<!-- TEST_BADGE --><img src="https://img.shields.io/badge/tests-passing-%2318963e?style=for-the-badge&logo=textpattern&logoColor=%23ffffff&logoSize=32&label=tests&labelColor=%23034015&color=%2318963e&cacheSeconds=600" alt="Project Tests Passing"><!-- END_TEST_BADGE -->
&nbsp;
<img src="https://img.shields.io/github/license/dev-kas/virtlang-go?style=for-the-badge&logo=2fas&logoColor=%23ffffff&logoSize=64&labelColor=%23701e25&color=%23ab3841&cacheSeconds=6000" alt="Project License">
&nbsp;
<img src="https://img.shields.io/github/stars/dev-kas/virtlang-go?style=for-the-badge&logo=github&logoSize=64&labelColor=%231d6791&color=%233894c9" alt="GitHub Repo stars">
&nbsp;
<img src="https://img.shields.io/github/v/release/dev-kas/virtlang-go?sort=semver&display_name=release&style=for-the-badge&logo=verizon&labelColor=%23347039&color=%234dab55&cacheSeconds=600" alt="GitHub Release">
</div>

# VirtLang
VirtLang is a custom-built programming language interpreter written in **Go**, designed for extensibility, clarity, and hackability. It powers a full-language runtime with native support for variables, functions, control flow, error handling, and complex types.

## 🚀 Overview

VirtLang features a clean three-stage architecture:

1. **Lexing** (`lexer.Tokenize()`) — converts source code into tokens  
2. **Parsing** (`parser.ProduceAST()`) — builds an Abstract Syntax Tree (AST)  
3. **Evaluation** (`evaluator.Evaluate()`) — runs the AST and produces results

The AST separates logic into:
- `Stmt` → non-returning statements (like `let`, `while`, `if`)
- `Expr` → return-producing expressions (like `1 + 2`, `"hello"`)

Runtime is powered by:
- A unified `RuntimeValue` type system
- Scoped environments for variable resolution
- Support for both **native** and **user-defined** functions

## 🧠 Language Features

- `let` and `const` declarations
- Functions (closures, parameters, call support)
- Control flow: `if`, `else`, `while`
- Structured error handling: `try`, `catch`
- Rich type system: numbers, strings, booleans, arrays, objects
- Member access (`obj.key`) and array indexing (`arr[i]`)
- Binary, logical, and comparison operators

## 🧪 Getting Started

Here's a minimal example showing how to evaluate VirtLang code in Go:

```go
package main

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/environment"
	"github.com/dev-kas/virtlang-go/evaluator"
	"github.com/dev-kas/virtlang-go/parser"
)

func main() {
	code := `let n = 1
while (n < 18) {
	n = n + 1
}
n // output: 18`

	// Set up parser and global environment
	p := parser.New()
	env := environment.NewEnvironment(nil)

	// Parse source code into AST
	program, err := p.ProduceAST(code)
	if err != nil {
		fmt.Printf("Syntax error: %v\n", err)
		return
	}

	// Evaluate AST
	result, runErr := evaluator.Evaluate(program, &env)
	if runErr != nil {
		fmt.Printf("Runtime error: %v\n", runErr)
		return
	}

	// Display result
	fmt.Printf("Result: %v (Type: %v)\n", result.Value, result.Type)
}
```

## 📚 Documentation

* Auto-generated Go package docs: [`DOCS.md`](DOCS.md)
* Full design write-up and architecture: [VirtLang Wiki](https://deepwiki.com/dev-kas/virtlang-go)

## 🤝 Contributing

Found a bug? Want to add a feature? See [`CONTRIBUTING.md`](CONTRIBUTING.md) for guidelines.

## 📊 Analytics

<div align="center">
  <img src="https://repobeats.axiom.co/api/embed/09a765e0d0bf50cf5dcc409272f31b3c66aa4b7c.svg" title="Repobeats analytics image for virtlang-go" alt="Repobeats analytics image for virtlang-go">
</div>
