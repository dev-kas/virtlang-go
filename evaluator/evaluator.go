package evaluator

import (
	"VirtLang/ast"
	"VirtLang/environment"
	"VirtLang/errors"
	"VirtLang/shared"
	"VirtLang/values"
	"fmt"
)

func Evaluate(astNode ast.Stmt, env *environment.Environment) (*shared.RuntimeValue, *errors.RuntimeError) {
	type_ := astNode.GetType()
	switch type_ {
	case ast.NumericLiteralNode:
		result := values.MK_NUMBER(int(astNode.(*ast.NumericLiteral).Value)) // TODO: convert to float64 later
		return &result, nil

	case ast.StringLiteralNode:
		result := values.MK_STRING(astNode.(*ast.StringLiteral).Value)
		return &result, nil

	case ast.IdentifierNode:
		return evalIdentifier(astNode.(*ast.Identifier), env)

	case ast.ObjectLiteralNode:
		return evalObjectExpr(astNode.(*ast.ObjectLiteral), env)

	case ast.CallExprNode:
		return evalCallExpr(astNode.(*ast.CallExpr), env)

	case ast.MemberExprNode:
		return evalMemberExpr(astNode.(*ast.MemberExpr), env)

	case ast.VarAssignmentExprNode:
		return evalVarAssignment(astNode.(*ast.VarAssignmentExpr), env)

	case ast.BinaryExprNode:
		return evalBinEx(astNode.(*ast.BinaryExpr), env)

	case ast.CompareExprNode:
		return evalComEx(astNode.(*ast.CompareExpr), env)

	case ast.ProgramNode:
		return evalProgram(astNode.(*ast.Program), env)

	case ast.VarDeclarationNode:
		return evalVarDecl(astNode.(*ast.VarDeclaration), env)

	case ast.FnDeclarationNode:
		return evalFnDecl(astNode.(*ast.FnDeclaration), env)

	case ast.IfStatementNode:
		return evalIfStmt(astNode.(*ast.IfStatement), env)

	case ast.WhileLoopNode:
		return evalWhileLoop(astNode.(*ast.WhileLoop), env)

	case ast.TryCatchStmtNode:
		return evalTryCatch(astNode.(*ast.TryCatchStmt), env)

	case ast.ReturnStmtNode:
		return evalReturnStmt(astNode.(*ast.ReturnStmt), env)

	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Unknown node type: %s", type_),
		}
	}
}
