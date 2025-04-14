package ast_test

import (
	"testing"

	"VirtLang/ast"
)

func TestNodeType_String(t *testing.T) {
	tests := []struct {
		nodeType ast.NodeType
		expected string
	}{
		{ast.ProgramNode, "Program"},
		{ast.VarDeclarationNode, "VarDeclaration"},
		{ast.FnDeclarationNode, "FnDeclaration"},
		{ast.IfStatementNode, "IfStatement"},
		{ast.WhileLoopNode, "WhileLoop"},
		{ast.VarAssignmentExprNode, "VarAssignmentExpr"},
		{ast.TryCatchStmtNode, "TryCatchStmt"},
		{ast.MemberExprNode, "MemberExpr"},
		{ast.ReturnStmtNode, "ReturnStmt"},
		{ast.CallExprNode, "CallExpr"},
		{ast.PropertyNode, "Property"},
		{ast.ObjectLiteralNode, "ObjectLiteral"},
		{ast.NumericLiteralNode, "NumericLiteral"},
		{ast.StringLiteralNode, "StringLiteral"},
		{ast.IdentifierNode, "Identifier"},
		{ast.CompareExprNode, "CompareExpr"},
		{ast.BinaryExprNode, "BinaryExpr"},
		{ast.NodeType(999), "UnknownNodeType"}, // Test the default case
	}

	for _, test := range tests {
		actual := test.nodeType.String()
		if actual != test.expected {
			t.Errorf("For NodeType %d, expected '%s', but got '%s'", test.nodeType, test.expected, actual)
		}
	}
}

func TestCompareOperatorConstants(t *testing.T) {
	if ast.LessThanOrEqual != "<=" {
		t.Errorf("Expected LessThanOrEqual to be '<=', but got '%s'", ast.LessThanOrEqual)
	}
	if ast.LessThan != "<" {
		t.Errorf("Expected LessThan to be '<', but got '%s'", ast.LessThan)
	}
	if ast.GreaterThanOrEqual != "=>" {
		t.Errorf("Expected GreaterThanOrEqual to be '=>', but got '%s'", ast.GreaterThanOrEqual)
	}
	if ast.GreaterThan != ">" {
		t.Errorf("Expected GreaterThan to be '>', but got '%s'", ast.GreaterThan)
	}
	if ast.Equal != "==" {
		t.Errorf("Expected Equal to be '==', but got '%s'", ast.Equal)
	}
	if ast.NotEqual != "!=" {
		t.Errorf("Expected NotEqual to be '!=', but got '%s'", ast.NotEqual)
	}
}

func TestBinaryOperatorConstants(t *testing.T) {
	if ast.Plus != "+" {
		t.Errorf("Expected Plus to be '+', but got '%s'", ast.Plus)
	}
	if ast.Minus != "-" {
		t.Errorf("Expected Minus to be '-', but got '%s'", ast.Minus)
	}
	if ast.Multiply != "*" {
		t.Errorf("Expected Multiply to be '*', but got '%s'", ast.Multiply)
	}
	if ast.Divide != "/" {
		t.Errorf("Expected Divide to be '/', but got '%s'", ast.Divide)
	}
	if ast.Modulo != "%" {
		t.Errorf("Expected Modulo to be '%%', but got '%s'", ast.Modulo)
	}
}

func TestStmtGetTypeMethods(t *testing.T) {
	stmtTests := []struct {
		name     string
		node     ast.Stmt
		expected ast.NodeType
	}{
		{"Program", &ast.Program{}, ast.ProgramNode},
		{"VarDeclaration", &ast.VarDeclaration{}, ast.VarDeclarationNode},
		{"FnDeclaration", &ast.FnDeclaration{}, ast.FnDeclarationNode},
		{"IfStatement", &ast.IfStatement{}, ast.IfStatementNode},
		{"WhileLoop", &ast.WhileLoop{}, ast.WhileLoopNode},
		{"TryCatchStmt", &ast.TryCatchStmt{}, ast.TryCatchStmtNode},
		{"ReturnStmt", &ast.ReturnStmt{}, ast.ReturnStmtNode},
	}

	for _, test := range stmtTests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.node.GetType()
			if actual != test.expected {
				t.Errorf("Expected %s.GetType() to return %v, but got %v", test.name, test.expected, actual)
			}
		})
	}
}

func TestExprGetTypeMethods(t *testing.T) {
	exprTests := []struct {
		name     string
		node     ast.Expr
		expected ast.NodeType
	}{
		{"VarAssignmentExpr", &ast.VarAssignmentExpr{}, ast.VarAssignmentExprNode},
		{"BinaryExpr", &ast.BinaryExpr{}, ast.BinaryExprNode},
		{"CompareExpr", &ast.CompareExpr{}, ast.CompareExprNode},
		{"CallExpr", &ast.CallExpr{}, ast.CallExprNode},
		{"MemberExpr", &ast.MemberExpr{}, ast.MemberExprNode},
		{"Identifier", &ast.Identifier{}, ast.IdentifierNode},
		{"NumericLiteral", &ast.NumericLiteral{}, ast.NumericLiteralNode},
		{"StringLiteral", &ast.StringLiteral{}, ast.StringLiteralNode},
		{"Property", &ast.Property{}, ast.PropertyNode},
		{"ObjectLiteral", &ast.ObjectLiteral{}, ast.ObjectLiteralNode},
	}

	for _, test := range exprTests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.node.GetType()
			if actual != test.expected {
				t.Errorf("Expected %s.GetType() to return %v, but got %v", test.name, test.expected, actual)
			}
		})
	}
}
