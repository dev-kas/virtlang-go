package evaluator

import (
	"fmt"

	"github.com/dev-kas/virtlang-go/v4/ast"
	"github.com/dev-kas/virtlang-go/v4/debugger"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/errors"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

// `dbgr` circulates the evaluator
func Evaluate(astNode ast.Stmt, env *environment.Environment, dbgr *debugger.Debugger) (*shared.RuntimeValue, *errors.RuntimeError) {
	type_ := astNode.GetType()

	// If debugger is attached
	if dbgr != nil {
		dbgr.CurrentFile = astNode.GetSourceMetadata().Filename
		dbgr.CurrentLine = astNode.GetSourceMetadata().StartLine
		if dbgr.IsDebuggable(type_) {
			if dbgr.ShouldStop(astNode.GetSourceMetadata().Filename, astNode.GetSourceMetadata().StartLine) {
				dbgr.Pause()
			}
			// Only wait if the node is debuggable
			// This is required to prevent debugging
			// at a microscopic level
			dbgr.WaitIfPaused(type_)
		}
	}

	switch type_ {
	case ast.NumericLiteralNode:
		result := values.MK_NUMBER(astNode.(*ast.NumericLiteral).Value)
		return &result, nil

	case ast.StringLiteralNode:
		result := values.MK_STRING(astNode.(*ast.StringLiteral).Value)
		return &result, nil

	case ast.IdentifierNode:
		return evalIdentifier(astNode.(*ast.Identifier), env)

	case ast.ObjectLiteralNode:
		return evalObjectExpr(astNode.(*ast.ObjectLiteral), env, dbgr)

	case ast.ArrayLiteralNode:
		return evalArrayExpr(astNode.(*ast.ArrayLiteral), env, dbgr)

	case ast.CallExprNode:
		return evalCallExpr(astNode.(*ast.CallExpr), env, dbgr)

	case ast.MemberExprNode:
		return evalMemberExpr(astNode.(*ast.MemberExpr), env, dbgr)

	case ast.VarAssignmentExprNode:
		return evalVarAssignment(astNode.(*ast.VarAssignmentExpr), env, dbgr)

	case ast.BinaryExprNode:
		return evalBinEx(astNode.(*ast.BinaryExpr), env, dbgr)

	case ast.CompareExprNode:
		return evalComEx(astNode.(*ast.CompareExpr), env, dbgr)

	case ast.LogicalExprNode:
		return evalLogicEx(astNode.(*ast.LogicalExpr), env, dbgr)

	case ast.ProgramNode:
		return evalProgram(astNode.(*ast.Program), env, dbgr)

	case ast.VarDeclarationNode:
		return evalVarDecl(astNode.(*ast.VarDeclaration), env, dbgr)

	case ast.FnDeclarationNode:
		return evalFnDecl(astNode.(*ast.FnDeclaration), env)

	case ast.IfStatementNode:
		return evalIfStmt(astNode.(*ast.IfStatement), env, dbgr)

	case ast.WhileLoopNode:
		return evalWhileLoop(astNode.(*ast.WhileLoop), env, dbgr)

	case ast.TryCatchStmtNode:
		return evalTryCatch(astNode.(*ast.TryCatchStmt), env, dbgr)

	case ast.ReturnStmtNode:
		return evalReturnStmt(astNode.(*ast.ReturnStmt), env, dbgr)

	case ast.BreakStmtNode:
		return evalBreakStmt(astNode.(*ast.BreakStmt), env)

	case ast.ContinueStmtNode:
		return evalContinueStmt(astNode.(*ast.ContinueStmt), env)

	case ast.ClassNode:
		return evalClass(astNode.(*ast.Class), env)

	case ast.ClassMethodNode:
		return evalClassMethod(astNode.(*ast.ClassMethod), env)

	case ast.ClassPropertyNode:
		return evalClassProperty(astNode.(*ast.ClassProperty), env, dbgr)

	case ast.DestructureDeclarationNode:
		return evalDestructureDeclaration(astNode.(*ast.DestructureDeclaration), env, dbgr)

	default:
		return nil, &errors.RuntimeError{
			Message: fmt.Sprintf("Unknown node type: %s", type_),
		}
	}
}
