package ast

type NodeType int

const (
	ProgramNode NodeType = iota
	VarDeclarationNode
	FnDeclarationNode
	IfStatementNode
	WhileLoopNode
	VarAssignmentExprNode
	TryCatchStmtNode
	MemberExprNode
	ReturnStmtNode
	CallExprNode
	PropertyNode
	ObjectLiteralNode
	NumericLiteralNode
	StringLiteralNode
	IdentifierNode
	CompareExprNode
	BinaryExprNode
)

func (n NodeType) String() string {
	switch n {
	case ProgramNode:
		return "Program"
	case VarDeclarationNode:
		return "VarDeclaration"
	case FnDeclarationNode:
		return "FnDeclaration"
	case IfStatementNode:
		return "IfStatement"
	case WhileLoopNode:
		return "WhileLoop"
	case VarAssignmentExprNode:
		return "VarAssignmentExpr"
	case TryCatchStmtNode:
		return "TryCatchStmt"
	case MemberExprNode:
		return "MemberExpr"
	case ReturnStmtNode:
		return "ReturnStmt"
	case CallExprNode:
		return "CallExpr"
	case PropertyNode:
		return "Property"
	case ObjectLiteralNode:
		return "ObjectLiteral"
	case NumericLiteralNode:
		return "NumericLiteral"
	case StringLiteralNode:
		return "StringLiteral"
	case IdentifierNode:
		return "Identifier"
	case CompareExprNode:
		return "CompareExpr"
	case BinaryExprNode:
		return "BinaryExpr"
	default:
		return "UnknownNodeType"
	}
}

type CompareOperator string

const (
	LessThanEqual    CompareOperator = "<="
	LessThan         CompareOperator = "<"
	GreaterThanEqual CompareOperator = "=>"
	GreaterThan      CompareOperator = ">"
	Equal            CompareOperator = "=="
	NotEqual         CompareOperator = "!="
)

type BinaryOperator string

const (
	Plus     BinaryOperator = "+"
	Minus    BinaryOperator = "-"
	Multiply BinaryOperator = "*"
	Divide   BinaryOperator = "/"
	Modulo   BinaryOperator = "%"
)

type Stmt interface {
	GetType() NodeType
}

type Expr interface {
	GetType() NodeType
}

// Statements

type Program struct {
	Stmts []Stmt
}

func (p *Program) GetType() NodeType { return ProgramNode }

type VarDeclaration struct {
	Constant   bool
	Identifier string
	Value      Expr
}

func (v *VarDeclaration) GetType() NodeType { return VarDeclarationNode }

type TryCatchStmt struct {
	Try      []Stmt
	Catch    []Stmt
	CatchVar string
}

func (t *TryCatchStmt) GetType() NodeType { return TryCatchStmtNode }

type FnDeclaration struct {
	Params    []string
	Name      string
	Body      []Stmt
	Async     bool
	Anonymous bool
}

func (f *FnDeclaration) GetType() NodeType { return FnDeclarationNode }

type IfStatement struct {
	Body      []Stmt
	Condition Expr
}

func (i *IfStatement) GetType() NodeType { return IfStatementNode }

type WhileLoop struct {
	Body      []Stmt
	Condition Expr
}

func (w *WhileLoop) GetType() NodeType { return WhileLoopNode }

type ReturnStmt struct {
	Value Expr
}

func (r *ReturnStmt) GetType() NodeType { return ReturnStmtNode }

// Expressions

type VarAssignmentExpr struct {
	Assignee Expr
	Value    Expr
}

func (v *VarAssignmentExpr) GetType() NodeType { return VarAssignmentExprNode }

type BinaryExpr struct {
	LHS      Expr
	RHS      Expr
	Operator BinaryOperator
}

func (b *BinaryExpr) GetType() NodeType { return BinaryExprNode }

type CompareExpr struct {
	LHS      Expr
	RHS      Expr
	Operator CompareOperator
}

func (c *CompareExpr) GetType() NodeType { return CompareExprNode }

type CallExpr struct {
	Args   []Expr
	Callee Expr
}

func (c *CallExpr) GetType() NodeType { return CallExprNode }

type MemberExpr struct {
	Object   Expr
	Value    Expr
	Computed bool
}

func (m *MemberExpr) GetType() NodeType { return MemberExprNode }

type Identifier struct {
	Symbol string
}

func (i *Identifier) GetType() NodeType { return IdentifierNode }

type NumericLiteral struct {
	Value int // TODO: Switch to float64
}

func (n *NumericLiteral) GetType() NodeType { return NumericLiteralNode }

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) GetType() NodeType { return StringLiteralNode }

type Property struct {
	Key   string
	Value Expr
}

func (p *Property) GetType() NodeType { return PropertyNode }

type ObjectLiteral struct {
	Properties []Property
}

func (o *ObjectLiteral) GetType() NodeType { return ObjectLiteralNode }
