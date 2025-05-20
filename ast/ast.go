package ast

type NodeType int

const (
	ProgramNode NodeType = iota
	VarDeclarationNode
	FnDeclarationNode
	IfStatementNode
	ClassNode
	ClassMethodNode
	ClassPropertyNode
	WhileLoopNode
	VarAssignmentExprNode
	TryCatchStmtNode
	MemberExprNode
	ReturnStmtNode
	BreakStmtNode
	ContinueStmtNode
	CallExprNode
	PropertyNode
	ObjectLiteralNode
	ArrayLiteralNode
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
	case ArrayLiteralNode:
		return "ArrayLiteral"
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
	case ClassNode:
		return "Class"
	case ClassMethodNode:
		return "ClassMethod"
	case ClassPropertyNode:
		return "ClassProperty"
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

type SourceMetadata struct {
	StartLine     int
	StartColumn   int
	EndLine       int
	EndColumn     int
	Filename      string
}

type Stmt interface {
	GetType() NodeType
}

type Expr interface {
	GetType() NodeType
}

// Statements

type Program struct {
	Stmts []Stmt
	SourceMetadata
}

func (p *Program) GetType() NodeType { return ProgramNode }

type VarDeclaration struct {
	Constant   bool
	Identifier string
	Value      Expr
	SourceMetadata
}

func (v *VarDeclaration) GetType() NodeType { return VarDeclarationNode }

type TryCatchStmt struct {
	Try      []Stmt
	Catch    []Stmt
	CatchVar string
	SourceMetadata
}

func (t *TryCatchStmt) GetType() NodeType { return TryCatchStmtNode }

type FnDeclaration struct {
	Params    []string
	Name      string
	Body      []Stmt
	Anonymous bool
	SourceMetadata
}

func (f *FnDeclaration) GetType() NodeType { return FnDeclarationNode }

type IfStatement struct {
	Body      []Stmt
	Condition Expr
	Else      []Stmt         // Optional else block
	ElseIf    []*IfStatement // Optional else-if branches
	SourceMetadata
}

func (i *IfStatement) GetType() NodeType { return IfStatementNode }

type Class struct {
	Name        string
	Body        []Stmt
	Constructor *ClassMethod
	SourceMetadata
}

func (c *Class) GetType() NodeType { return ClassNode }

type ClassMethod struct {
	Name     string
	Body     []Stmt
	Params   []string
	IsPublic bool
	SourceMetadata
}

func (c *ClassMethod) GetType() NodeType { return ClassMethodNode }

type ClassProperty struct {
	Name     string
	Value    Expr
	IsPublic bool
	SourceMetadata
}

func (c *ClassProperty) GetType() NodeType { return ClassPropertyNode }

type WhileLoop struct {
	Body      []Stmt
	Condition Expr
	SourceMetadata
}

func (w *WhileLoop) GetType() NodeType { return WhileLoopNode }

type ReturnStmt struct {
	Value Expr
	SourceMetadata
}

func (r *ReturnStmt) GetType() NodeType { return ReturnStmtNode }

type BreakStmt struct {
	SourceMetadata
}

func (r *BreakStmt) GetType() NodeType { return BreakStmtNode }

type ContinueStmt struct {
	SourceMetadata
}

func (r *ContinueStmt) GetType() NodeType { return ContinueStmtNode }

// Expressions

type VarAssignmentExpr struct {
	Assignee Expr
	Value    Expr
	SourceMetadata
}

func (v *VarAssignmentExpr) GetType() NodeType { return VarAssignmentExprNode }

type BinaryExpr struct {
	LHS      Expr
	RHS      Expr
	Operator BinaryOperator
	SourceMetadata
}

func (b *BinaryExpr) GetType() NodeType { return BinaryExprNode }

type CompareExpr struct {
	LHS      Expr
	RHS      Expr
	Operator CompareOperator
	SourceMetadata
}

func (c *CompareExpr) GetType() NodeType { return CompareExprNode }

type CallExpr struct {
	Args   []Expr
	Callee Expr
	SourceMetadata
}

func (c *CallExpr) GetType() NodeType { return CallExprNode }

type MemberExpr struct {
	Object   Expr
	Value    Expr
	Computed bool
	SourceMetadata
}

func (m *MemberExpr) GetType() NodeType { return MemberExprNode }

type Identifier struct {
	Symbol string
	SourceMetadata
}

func (i *Identifier) GetType() NodeType { return IdentifierNode }

type NumericLiteral struct {
	Value float64
	SourceMetadata
}

func (n *NumericLiteral) GetType() NodeType { return NumericLiteralNode }

type StringLiteral struct {
	Value string
	SourceMetadata
}

func (s *StringLiteral) GetType() NodeType { return StringLiteralNode }

type Property struct {
	Key   string
	Value Expr
	SourceMetadata
}

func (p *Property) GetType() NodeType { return PropertyNode }

type ObjectLiteral struct {
	Properties []Property
	SourceMetadata
}

func (o *ObjectLiteral) GetType() NodeType { return ObjectLiteralNode }

type ArrayLiteral struct {
	Elements []Expr
	SourceMetadata
}

func (o *ArrayLiteral) GetType() NodeType { return ArrayLiteralNode }
