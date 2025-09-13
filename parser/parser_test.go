package parser_test

import (
	"testing"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/internal/testhelpers"
	"github.com/dev-kas/virtlang-go/v4/parser"
)

func TestBinaryExpr(t *testing.T) {
	srccode := "1 + 2 * (3 - 4) / 5"
	p := parser.New("test")
	prog, err := p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if prog.GetType() != ast.ProgramNode {
		t.Fatalf("Didn't receive a program, got %s", prog.GetType())
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	add := prog.Stmts[0]
	if add.GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected top-level to be BinaryExpr (+), got %s", add.GetType())
	}
	if add.(*ast.BinaryExpr).Operator != ast.Plus {
		t.Fatalf("Expected operator to be Plus (+), got %s", add.(*ast.BinaryExpr).Operator)
	}

	left := add.(*ast.BinaryExpr).LHS
	if left.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (+) to be NumericLiteral, got %s", left.GetType())
	}
	if left.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected left side of (+) to be 1, got %f", left.(*ast.NumericLiteral).Value)
	}

	divide := add.(*ast.BinaryExpr).RHS
	if divide.GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected right side of (+) to be BinaryExpr (/), got %s", divide.GetType())
	}
	if divide.(*ast.BinaryExpr).Operator != ast.Divide {
		t.Fatalf("Expected operator to be Divide (/), got %s", divide.(*ast.BinaryExpr).Operator)
	}

	divideRight := divide.(*ast.BinaryExpr).RHS
	if divideRight.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (/) to be NumericLiteral, got %s", divideRight.GetType())
	}
	if divideRight.(*ast.NumericLiteral).Value != 5 {
		t.Fatalf("Expected right side of (/) to be 5, got %f", divideRight.(*ast.NumericLiteral).Value)
	}

	multiply := divide.(*ast.BinaryExpr).LHS
	if multiply.GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected left side of (/) to be BinaryExpr (*), got %s", multiply.GetType())
	}
	if multiply.(*ast.BinaryExpr).Operator != ast.Multiply {
		t.Fatalf("Expected operator to be Multiply (*), got %s", multiply.(*ast.BinaryExpr).Operator)
	}

	multiplyLeft := multiply.(*ast.BinaryExpr).LHS
	if multiplyLeft.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (*) to be NumericLiteral, got %s", multiplyLeft.GetType())
	}
	if multiplyLeft.(*ast.NumericLiteral).Value != 2 {
		t.Fatalf("Expected left side of (*) to be 2, got %f", multiplyLeft.(*ast.NumericLiteral).Value)
	}

	subtract := multiply.(*ast.BinaryExpr).RHS
	if subtract.GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected right side of (*) to be BinaryExpr (-), got %s", subtract.GetType())
	}
	if subtract.(*ast.BinaryExpr).Operator != ast.Minus {
		t.Fatalf("Expected operator to be Minus (-), got %s", subtract.(*ast.BinaryExpr).Operator)
	}

	subtractLeft := subtract.(*ast.BinaryExpr).LHS
	if subtractLeft.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (-) to be NumericLiteral, got %s", subtractLeft.GetType())
	}
	if subtractLeft.(*ast.NumericLiteral).Value != 3 {
		t.Fatalf("Expected left side of (-) to be 3, got %f", subtractLeft.(*ast.NumericLiteral).Value)
	}

	subtractRight := subtract.(*ast.BinaryExpr).RHS
	if subtractRight.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (-) to be NumericLiteral, got %s", subtractRight.GetType())
	}
	if subtractRight.(*ast.NumericLiteral).Value != 4 {
		t.Fatalf("Expected right side of (-) to be 4, got %f", subtractRight.(*ast.NumericLiteral).Value)
	}

	// commonly made mistakes

	tests := []string{
		"1 +",
		"1 -",
		"1 *",
		"1 /",
		"1 %",
		"+ 1",
		"* 1",
		"/ 1",
		"% 1",
	}

	for _, test := range tests {
		_, err := p.ProduceAST(test)
		if err == nil {
			t.Fatalf("Expected error for %s", test)
		}
	}
}

func TestTryCatchExpr(t *testing.T) {
	srccode := `try {
	2 + 1
} catch e {
	4 - 1
	}`

	p := parser.New("test")

	prog, err := p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt_tc := prog.Stmts[0]

	if stmt_tc.GetType() != ast.TryCatchStmtNode {
		t.Fatalf("Expected a TryCatchNode, got %s", stmt_tc.GetType())
	}

	tblock := stmt_tc.(*ast.TryCatchStmt).Try
	if len(tblock) != 1 {
		t.Fatalf("Expected 1 statement in try block, got %d", len(tblock))
	}

	add := tblock[0]
	if add.GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected a BinaryExprNode, got %s", add.GetType())
	}

	lhs := add.(*ast.BinaryExpr).LHS
	if lhs.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (+) to be NumericLiteral, got %s", lhs.GetType())
	}
	if lhs.(*ast.NumericLiteral).Value != 2 {
		t.Fatalf("Expected left side of (+) to be 1, got %f", lhs.(*ast.NumericLiteral).Value)
	}

	rhs := add.(*ast.BinaryExpr).RHS
	if rhs.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (+) to be NumericLiteral, got %s", rhs.GetType())
	}
	if rhs.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected right side of (+) to be 2, got %f", rhs.(*ast.NumericLiteral).Value)
	}

	catchVar := stmt_tc.(*ast.TryCatchStmt).CatchVar
	if catchVar != "e" {
		t.Fatalf("Expected catch variable to be e, got %s", catchVar)
	}

	catch := stmt_tc.(*ast.TryCatchStmt).Catch
	if len(catch) != 1 {
		t.Fatalf("Expected 1 statement in catch block, got %d", len(catch))
	}

	subtract := catch[0]
	if subtract.GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected a BinaryExprNode, got %s", subtract.GetType())
	}

	lhs = subtract.(*ast.BinaryExpr).LHS
	if lhs.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (-) to be NumericLiteral, got %s", lhs.GetType())
	}
	if lhs.(*ast.NumericLiteral).Value != 4 {
		t.Fatalf("Expected left side of (-) to be 4, got %f", lhs.(*ast.NumericLiteral).Value)
	}

	rhs = subtract.(*ast.BinaryExpr).RHS
	if rhs.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (-) to be NumericLiteral, got %s", rhs.GetType())
	}
	if rhs.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected right side of (-) to be 1, got %f", rhs.(*ast.NumericLiteral).Value)
	}

	// commonly made mistakes

	tests := []string{
		`try {a+b} catch {a-b}`,
		`try {a+b} e {a-b}`,
		`try {a+b} catch e`,
		`try a+b catch e {a-b}`,
		`try {a+b} catch e a-b`,
		`try a+b catch e a-b`,
	}

	for _, test := range tests {
		_, err = p.ProduceAST(test)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	}
}

func TestVarStuff(t *testing.T) {
	srccode := "let a = 1"

	p := parser.New("test")

	prog, err := p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt := prog.Stmts[0]

	if stmt.GetType() != ast.VarDeclarationNode {
		t.Fatalf("Expected a VarAssignmentNode, got %s", stmt.GetType())
	}

	if stmt.(*ast.VarDeclaration).Constant {
		t.Fatalf("Expected constant to be false, got %t", stmt.(*ast.VarDeclaration).Constant)
	}

	if stmt.(*ast.VarDeclaration).Identifier != "a" {
		t.Fatalf("Expected identifier to be a, got %s", stmt.(*ast.VarDeclaration).Identifier)
	}

	if stmt.(*ast.VarDeclaration).Value.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected value to be NumericLiteral, got %s", stmt.(*ast.VarDeclaration).Value.GetType())
	}

	if stmt.(*ast.VarDeclaration).Value.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected value to be 1, got %f", stmt.(*ast.VarDeclaration).Value.(*ast.NumericLiteral).Value)
	}

	srccode = "const b = a + 3"

	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt = prog.Stmts[0]

	if stmt.GetType() != ast.VarDeclarationNode {
		t.Fatalf("Expected a VarAssignmentNode, got %s", stmt.GetType())
	}

	if !stmt.(*ast.VarDeclaration).Constant {
		t.Fatalf("Expected constant to be true, got %t", stmt.(*ast.VarDeclaration).Constant)
	}

	if stmt.(*ast.VarDeclaration).Identifier != "b" {
		t.Fatalf("Expected identifier to be b, got %s", stmt.(*ast.VarDeclaration).Identifier)
	}

	if stmt.(*ast.VarDeclaration).Value.GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected value to be BinaryExpr, got %s", stmt.(*ast.VarDeclaration).Value.GetType())
	}

	if stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).Operator != ast.Plus {
		t.Fatalf("Expected operator to be +, got %s", stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).Operator)
	}

	if stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).LHS.GetType() != ast.IdentifierNode {
		t.Fatalf("Expected left side of (+) to be Identifier, got %s", stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).LHS.GetType())
	}

	if stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).LHS.(*ast.Identifier).Symbol != "a" {
		t.Fatalf("Expected left side of (+) to be a, got %s", stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).LHS.(*ast.Identifier).Symbol)
	}

	if stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).RHS.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (+) to be NumericLiteral, got %s", stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).RHS.GetType())
	}

	if stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).RHS.(*ast.NumericLiteral).Value != 3 {
		t.Fatalf("Expected right side of (+) to be 3, got %f", stmt.(*ast.VarDeclaration).Value.(*ast.BinaryExpr).RHS.(*ast.NumericLiteral).Value)
	}

	srccode = "b = a - 3"

	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt = prog.Stmts[0]

	if stmt.GetType() != ast.VarAssignmentExprNode {
		t.Fatalf("Expected a VarAssignmentNode, got %s", stmt.GetType())
	}

	if stmt.(*ast.VarAssignmentExpr).Assignee.(*ast.Identifier).Symbol != "b" {
		t.Fatalf("Expected identifier to be b, got %s", stmt.(*ast.VarAssignmentExpr).Assignee.(*ast.Identifier).Symbol)
	}

	if stmt.(*ast.VarAssignmentExpr).Value.GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected value to be BinaryExpr, got %s", stmt.(*ast.VarAssignmentExpr).Value.GetType())
	}

	if stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).Operator != ast.Minus {
		t.Fatalf("Expected operator to be -, got %s", stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).Operator)
	}

	if stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).LHS.GetType() != ast.IdentifierNode {
		t.Fatalf("Expected left side of (-) to be Identifier, got %s", stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).LHS.GetType())
	}

	if stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).LHS.(*ast.Identifier).Symbol != "a" {
		t.Fatalf("Expected left side of (-) to be a, got %s", stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).LHS.(*ast.Identifier).Symbol)
	}

	if stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).RHS.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (-) to be NumericLiteral, got %s", stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).RHS.GetType())
	}

	if stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).RHS.(*ast.NumericLiteral).Value != 3 {
		t.Fatalf("Expected right side of (-) to be 3, got %f", stmt.(*ast.VarAssignmentExpr).Value.(*ast.BinaryExpr).RHS.(*ast.NumericLiteral).Value)
	}

	// Commonly made mistakes

	tests := []string{
		"let a",
		"let a =",
		"let = 3",
		"let =",
		"const a",
		"const a =",
		"const = 3",
		"const =",
	}

	for _, test := range tests {
		_, err := p.ProduceAST(test)
		if err == nil {
			t.Fatalf("Expected error for %q, got nil", test)
		}
	}
}

func TestComparison(t *testing.T) {
	srccode := "1 < 2"
	p := parser.New("test")
	prog, err := p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt := prog.Stmts[0]

	if stmt.GetType() != ast.CompareExprNode {
		t.Fatalf("Expected a CompareNode, got %s", stmt.GetType())
	}

	left := stmt.(*ast.CompareExpr).LHS
	if left.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (<) to be NumericLiteral, got %s", left.GetType())
	}
	if left.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected left side of (<) to be 1, got %f", left.(*ast.NumericLiteral).Value)
	}

	right := stmt.(*ast.CompareExpr).RHS
	if right.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (<) to be NumericLiteral, got %s", right.GetType())
	}
	if right.(*ast.NumericLiteral).Value != 2 {
		t.Fatalf("Expected right side of (<) to be 2, got %f", right.(*ast.NumericLiteral).Value)
	}

	if stmt.(*ast.CompareExpr).Operator != ast.LessThan {
		t.Fatalf("Expected operator to be <, got %s", stmt.(*ast.CompareExpr).Operator)
	}

	srccode = "1 > 2"
	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt = prog.Stmts[0]

	if stmt.GetType() != ast.CompareExprNode {
		t.Fatalf("Expected a CompareNode, got %s", stmt.GetType())
	}

	left = stmt.(*ast.CompareExpr).LHS
	if left.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (<) to be NumericLiteral, got %s", left.GetType())
	}
	if left.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected left side of (<) to be 1, got %f", left.(*ast.NumericLiteral).Value)
	}

	right = stmt.(*ast.CompareExpr).RHS
	if right.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (<) to be NumericLiteral, got %s", right.GetType())
	}
	if right.(*ast.NumericLiteral).Value != 2 {
		t.Fatalf("Expected right side of (<) to be 2, got %f", right.(*ast.NumericLiteral).Value)
	}

	if stmt.(*ast.CompareExpr).Operator != ast.GreaterThan {
		t.Fatalf("Expected operator to be >, got %s", stmt.(*ast.CompareExpr).Operator)
	}

	srccode = "1 == 2"
	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt = prog.Stmts[0]

	if stmt.GetType() != ast.CompareExprNode {
		t.Fatalf("Expected a CompareNode, got %s", stmt.GetType())
	}

	left = stmt.(*ast.CompareExpr).LHS
	if left.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (<) to be NumericLiteral, got %s", left.GetType())
	}
	if left.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected left side of (<) to be 1, got %f", left.(*ast.NumericLiteral).Value)
	}

	right = stmt.(*ast.CompareExpr).RHS
	if right.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (<) to be NumericLiteral, got %s", right.GetType())
	}
	if right.(*ast.NumericLiteral).Value != 2 {
		t.Fatalf("Expected right side of (<) to be 2, got %f", right.(*ast.NumericLiteral).Value)
	}

	if stmt.(*ast.CompareExpr).Operator != ast.Equal {
		t.Fatalf("Expected operator to be ==, got %s", stmt.(*ast.CompareExpr).Operator)
	}

	srccode = "1 != 2"
	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt = prog.Stmts[0]

	if stmt.GetType() != ast.CompareExprNode {
		t.Fatalf("Expected a CompareNode, got %s", stmt.GetType())
	}

	left = stmt.(*ast.CompareExpr).LHS
	if left.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (<) to be NumericLiteral, got %s", left.GetType())
	}
	if left.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected left side of (<) to be 1, got %f", left.(*ast.NumericLiteral).Value)
	}

	right = stmt.(*ast.CompareExpr).RHS
	if right.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (<) to be NumericLiteral, got %s", right.GetType())
	}
	if right.(*ast.NumericLiteral).Value != 2 {
		t.Fatalf("Expected right side of (<) to be 2, got %f", right.(*ast.NumericLiteral).Value)
	}

	if stmt.(*ast.CompareExpr).Operator != ast.NotEqual {
		t.Fatalf("Expected operator to be !=, got %s", stmt.(*ast.CompareExpr).Operator)
	}

	srccode = "1 <= 2"
	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt = prog.Stmts[0]

	if stmt.GetType() != ast.CompareExprNode {
		t.Fatalf("Expected a CompareNode, got %s", stmt.GetType())
	}

	left = stmt.(*ast.CompareExpr).LHS
	if left.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (<) to be NumericLiteral, got %s", left.GetType())
	}
	if left.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected left side of (<) to be 1, got %f", left.(*ast.NumericLiteral).Value)
	}

	right = stmt.(*ast.CompareExpr).RHS
	if right.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (<) to be NumericLiteral, got %s", right.GetType())
	}
	if right.(*ast.NumericLiteral).Value != 2 {
		t.Fatalf("Expected right side of (<) to be 2, got %f", right.(*ast.NumericLiteral).Value)
	}

	if stmt.(*ast.CompareExpr).Operator != ast.LessThanEqual {
		t.Fatalf("Expected operator to be <=, got %s", stmt.(*ast.CompareExpr).Operator)
	}

	srccode = "1 >= 2"
	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt = prog.Stmts[0]

	if stmt.GetType() != ast.CompareExprNode {
		t.Fatalf("Expected a CompareNode, got %s", stmt.GetType())
	}

	left = stmt.(*ast.CompareExpr).LHS
	if left.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected left side of (<) to be NumericLiteral, got %s", left.GetType())
	}
	if left.(*ast.NumericLiteral).Value != 1 {
		t.Fatalf("Expected left side of (<) to be 1, got %f", left.(*ast.NumericLiteral).Value)
	}

	right = stmt.(*ast.CompareExpr).RHS
	if right.GetType() != ast.NumericLiteralNode {
		t.Fatalf("Expected right side of (<) to be NumericLiteral, got %s", right.GetType())
	}
	if right.(*ast.NumericLiteral).Value != 2 {
		t.Fatalf("Expected right side of (<) to be 2, got %f", right.(*ast.NumericLiteral).Value)
	}

	if stmt.(*ast.CompareExpr).Operator != ast.GreaterThanEqual {
		t.Fatalf("Expected operator to be >=, got %s", stmt.(*ast.CompareExpr).Operator)
	}

	// complex

	srccode = "(1 < 2) == true"
	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt = prog.Stmts[0]

	if stmt.GetType() != ast.CompareExprNode {
		t.Fatalf("Expected a CompareNode, got %s", stmt.GetType())
	}

	left = stmt.(*ast.CompareExpr).LHS
	if left.GetType() != ast.CompareExprNode {
		t.Fatalf("Expected left side of (<) to be CompareExpr, got %s", left.GetType())
	}

	if left.(*ast.CompareExpr).Operator != ast.LessThan {
		t.Fatalf("Expected operator to be <, got %s", left.(*ast.BinaryExpr).Operator)
	}

	right = stmt.(*ast.CompareExpr).RHS
	if right.GetType() != ast.IdentifierNode {
		t.Fatalf("Expected right side of (<) to be Identifier, got %s", right.GetType())
	}
	if right.(*ast.Identifier).Symbol != "true" {
		t.Fatalf("Expected right side of (<) to be true, got %s", right.(*ast.Identifier).Symbol)
	}

	if stmt.(*ast.CompareExpr).Operator != ast.Equal {
		t.Fatalf("Expected operator to be ==, got %s", stmt.(*ast.CompareExpr).Operator)
	}

	// Commonly made mistakes

	tests := []string{
		"1 <",
		"1 >",
		"< 2",
		"> 2",
		"<= 2",
		">= 2",
	}

	for _, test := range tests {
		_, err := p.ProduceAST(test)
		if err == nil {
			t.Fatalf("Expected error for %q, got nil", test)
		}
	}
}

func TestFnDecl(t *testing.T) {
	srccode := `fn add(a, b) {
		a + b
	}`

	p := parser.New("test")
	prog, err := p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt := prog.Stmts[0]

	if stmt.GetType() != ast.FnDeclarationNode {
		t.Fatalf("Expected a FnDeclarationNode, got %s", stmt.GetType())
	}

	if stmt.(*ast.FnDeclaration).Name != "add" {
		t.Fatalf("Expected name to be add, got %s", stmt.(*ast.FnDeclaration).Name)
	}

	if len(stmt.(*ast.FnDeclaration).Params) != 2 {
		t.Fatalf("Expected 2 params, got %d", len(stmt.(*ast.FnDeclaration).Params))
	}

	if stmt.(*ast.FnDeclaration).Params[0] != "a" {
		t.Fatalf("Expected first param to be a, got %s", stmt.(*ast.FnDeclaration).Params[0])
	}

	if stmt.(*ast.FnDeclaration).Params[1] != "b" {
		t.Fatalf("Expected second param to be b, got %s", stmt.(*ast.FnDeclaration).Params[1])
	}

	if len(stmt.(*ast.FnDeclaration).Body) != 1 {
		t.Fatalf("Expected 1 statement in body, got %d", len(stmt.(*ast.FnDeclaration).Body))
	}

	if stmt.(*ast.FnDeclaration).Body[0].GetType() != ast.BinaryExprNode {
		t.Fatalf("Expected body to be BinaryExpr, got %s", stmt.(*ast.FnDeclaration).Body[0].GetType())
	}

	if stmt.(*ast.FnDeclaration).Body[0].(*ast.BinaryExpr).Operator != ast.Plus {
		t.Fatalf("Expected operator to be +, got %s", stmt.(*ast.FnDeclaration).Body[0].(*ast.BinaryExpr).Operator)
	}

	if stmt.(*ast.FnDeclaration).Body[0].(*ast.BinaryExpr).LHS.(*ast.Identifier).Symbol != "a" {
		t.Fatalf("Expected left side of (+) to be a, got %s", stmt.(*ast.FnDeclaration).Body[0].(*ast.BinaryExpr).LHS.(*ast.Identifier).Symbol)
	}

	if stmt.(*ast.FnDeclaration).Body[0].(*ast.BinaryExpr).RHS.(*ast.Identifier).Symbol != "b" {
		t.Fatalf("Expected right side of (+) to be b, got %s", stmt.(*ast.FnDeclaration).Body[0].(*ast.BinaryExpr).RHS.(*ast.Identifier).Symbol)
	}

	if stmt.(*ast.FnDeclaration).Anonymous {
		t.Fatalf("Expected anonymous to be false, got true")
	}

	// Commonly made mistakes

	tests := []string{
		"fn",
		"fn add",
		"fn add()",
		"fn add(a",
		"fn add(a,",
		"fn add(a, b",
		"fn add(a, b) {",
		"fn add(a, b) { a",
		"fn add(a, b) { a +",
		"fn add(a, b) { a + b",
		"fn { a + b }",
	}

	for _, test := range tests {
		_, err := p.ProduceAST(test)
		if err == nil {
			t.Fatalf("Expected error for %q, got nil", test)
		}
	}
}

func TestIfStmt(t *testing.T) {
	// Test basic if statement
	srccode := `if (3+1 > 3) {myFunction()}`
	p := parser.New("test")
	prog, err := p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	stmt := prog.Stmts[0]

	if stmt.GetType() != ast.IfStatementNode {
		t.Fatalf("Expected a IfStatementNode, got %s", stmt.GetType())
	}

	// Test if-else statement
	srccode = `if (3+1 > 3) {myFunction()} else {otherFunction()}`
	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	stmt = prog.Stmts[0]
	if stmt.(*ast.IfStatement).Else == nil {
		t.Fatalf("Expected an else block, got nil")
	}

	if len(stmt.(*ast.IfStatement).Else) == 0 {
		t.Fatalf("Expected else block to have statements, but it is empty")
	}

	if stmt.(*ast.IfStatement).Else[0].(*ast.CallExpr).Callee.(*ast.Identifier).Symbol != "otherFunction" {
		t.Fatalf("Expected else block to call otherFunction, got %s", stmt.(*ast.IfStatement).Else[0].(*ast.CallExpr).Callee.(*ast.Identifier).Symbol)
	}

	// Test if-else if-else chain
	srccode = `if (3+1 > 3) {myFunction()} else if (2 < 5) {anotherFunction()} else {finalFunction()}`
	prog, err = p.ProduceAST(srccode)
	if err != nil {
		t.Fatal(err)
	}

	stmt = prog.Stmts[0]
	if len(stmt.(*ast.IfStatement).ElseIf) != 1 {
		t.Fatalf("Expected 1 else-if block, got %d", len(stmt.(*ast.IfStatement).ElseIf))
	}

	if len(stmt.(*ast.IfStatement).ElseIf) > 0 {
		if stmt.(*ast.IfStatement).ElseIf[0].Condition.(*ast.CompareExpr).Operator != ast.LessThan {
			t.Fatalf("Expected else-if condition operator to be <, got %s", stmt.(*ast.IfStatement).ElseIf[0].Condition.(*ast.CompareExpr).Operator)
		}

		if len(stmt.(*ast.IfStatement).ElseIf[0].Body) > 0 {
			if stmt.(*ast.IfStatement).ElseIf[0].Body[0].(*ast.CallExpr).Callee.(*ast.Identifier).Symbol != "anotherFunction" {
				t.Fatalf("Expected else-if block to call anotherFunction, got %s", stmt.(*ast.IfStatement).ElseIf[0].Body[0].(*ast.CallExpr).Callee.(*ast.Identifier).Symbol)
			}
		}
	}

	if len(stmt.(*ast.IfStatement).Else) > 0 {
		if stmt.(*ast.IfStatement).Else[0].(*ast.CallExpr).Callee.(*ast.Identifier).Symbol != "finalFunction" {
			t.Fatalf("Expected else block to call finalFunction, got %s", stmt.(*ast.IfStatement).Else[0].(*ast.CallExpr).Callee.(*ast.Identifier).Symbol)
		}
	}

	// Commonly made mistakes

	tests := []string{
		"if (3+1 > 3",
		"if (3+1 > 3)",
		"if (3+1 > 3) { myFunction()",
		"if (3+1 > 3) myFunction()",
		"if 3+1 > 3 { myFunction() }",
		"if 3+1 > 3 myFunction()",
	}

	for _, test := range tests {
		_, err := p.ProduceAST(test)
		if err == nil {
			t.Fatalf("Expected error for %q, got nil", test)
		}
	}
}

// #7 Feature: Implement Destructuring Assignments (Object and Array)
func TestDestructurePatterns(t *testing.T) {
	tests := []struct {
		src           string
		expectType    string // "object" or "array"
		expectKeys    []string
		expectElCount int
	}{
		// object patterns
		{
			src:        `let { name, age } = { name: "Jedi", age: 30 }`,
			expectType: "object",
			expectKeys: []string{"name", "age"},
		},
		{
			src:        `let { user: { first, last }, id } = obj`,
			expectType: "object",
			expectKeys: []string{"user", "id"},
		},
		{
			src:        `let { x, y, ...rest } = point`,
			expectType: "object",
			expectKeys: []string{"x", "y"},
		},

		// array patterns
		{
			src:           `let [a, b, c] = [1,2,3]`,
			expectType:    "array",
			expectElCount: 3,
		},
		{
			src:           `let [first, , third] = arr`,
			expectType:    "array",
			expectElCount: 3,
		},
		{
			src:           `let [x, y, ...rest] = nums`,
			expectType:    "array",
			expectElCount: 2,
		},
		{
			src:           `let [[a,b], c] = nested`,
			expectType:    "array",
			expectElCount: 2,
		},
	}

	for _, tt := range tests {
		prog := testhelpers.MustParse(t, tt.src)
		if len(prog.Stmts) != 1 {
			t.Fatalf("[%s] Expected 1 statement, got %d", tt.src, len(prog.Stmts))
		}

		decl, ok := prog.Stmts[0].(*ast.DestructureDeclaration)
		if !ok {
			t.Fatalf("[%s] Expected DestructureDeclaration, got %s", tt.src, prog.Stmts[0].GetType())
		}

		switch tt.expectType {
		case "object":
			pattern, ok := decl.Pattern.(*ast.DestructureObjectPattern)
			if !ok {
				t.Fatalf("[%s] Expected ObjectPattern, got %s", tt.src, decl.Pattern.GetType())
			}
			if len(pattern.Properties) != len(tt.expectKeys) {
				t.Fatalf("[%s] Expected %d properties, got %d", tt.src, len(tt.expectKeys), len(pattern.Properties))
			}
			for i, key := range tt.expectKeys {
				if pattern.Properties[i].Key != key {
					t.Fatalf("[%s] Expected key %s, got %s", tt.src, key, pattern.Properties[i].Key)
				}
			}

		case "array":
			pattern, ok := decl.Pattern.(*ast.DestructureArrayPattern)
			if !ok {
				t.Fatalf("[%s] Expected ArrayPattern, got %s", tt.src, decl.Pattern.GetType())
			}
			if len(pattern.Elements) != tt.expectElCount {
				t.Fatalf("[%s] Expected %d elements, got %d", tt.src, tt.expectElCount, len(pattern.Elements))
			}
		}
	}

	// commonly made errors
	// objects
	testhelpers.ExpectParseError(t, "let { a, b, ...rest, c } = obj")        // rest not last
	testhelpers.ExpectParseError(t, "let { ...rest, a } = obj")              // rest not last
	testhelpers.ExpectParseError(t, "let { a, , c, d = 5 } = obj")           // skipping not allowed in object destructuring
	testhelpers.ExpectParseError(t, "let { a, b, ...rest, ...extra } = obj") // multiple rest

	// arrays
	testhelpers.ExpectParseError(t, "let [x, y, ...rest, z] = arr")     // rest not last
	testhelpers.ExpectParseError(t, "let [...rest, x] = arr")           // rest not last
	testhelpers.ExpectParseError(t, "let [a, ...rest, ...extra] = arr") // multiple rest
	testhelpers.ExpectParseError(t, "let [a, b, , ...rest, c] = arr")   // rest not last even with skipped element
}
