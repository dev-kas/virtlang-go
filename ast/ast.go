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
	LogicalExprNode
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
	case LogicalExprNode:
		return "LogicalExpr"
	default:
		return "UnknownNodeType"
	}
}

type CompareOperator string

const (
	LessThanEqual    CompareOperator = "<="
	LessThan         CompareOperator = "<"
	GreaterThanEqual CompareOperator = ">="
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

type LogicalOperator string

const (
	LogicalAND LogicalOperator = "&&"
	LogicalOR  LogicalOperator = "||"
	LogicalNilCoalescing LogicalOperator = "??"
	LogicalNOT LogicalOperator = "!"
)

type SourceMetadata struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Filename    string
}

type Stmt interface {
	GetType() NodeType
	GetSourceMetadata() SourceMetadata
}

type Expr interface {
	GetType() NodeType
	GetSourceMetadata() SourceMetadata
}

// Statements

type Program struct {
	Stmts []Stmt
	SourceMetadata
}

func (p *Program) GetType() NodeType                 { return ProgramNode }
func (p *Program) GetSourceMetadata() SourceMetadata { return p.SourceMetadata }

type VarDeclaration struct {
	Constant   bool
	Identifier string
	Value      Expr
	SourceMetadata
}

func (v *VarDeclaration) GetType() NodeType                 { return VarDeclarationNode }
func (v *VarDeclaration) GetSourceMetadata() SourceMetadata { return v.SourceMetadata }

type TryCatchStmt struct {
	Try      []Stmt
	Catch    []Stmt
	CatchVar string
	SourceMetadata
}

func (t *TryCatchStmt) GetType() NodeType                 { return TryCatchStmtNode }
func (t *TryCatchStmt) GetSourceMetadata() SourceMetadata { return t.SourceMetadata }

type FnDeclaration struct {
	Params    []string
	Name      string
	Body      []Stmt
	Anonymous bool
	SourceMetadata
}

func (f *FnDeclaration) GetType() NodeType                 { return FnDeclarationNode }
func (f *FnDeclaration) GetSourceMetadata() SourceMetadata { return f.SourceMetadata }

type IfStatement struct {
	Body      []Stmt
	Condition Expr
	Else      []Stmt         // Optional else block
	ElseIf    []*IfStatement // Optional else-if branches
	SourceMetadata
}

func (i *IfStatement) GetType() NodeType                 { return IfStatementNode }
func (i *IfStatement) GetSourceMetadata() SourceMetadata { return i.SourceMetadata }

type Class struct {
	Name        string
	Body        []Stmt
	Constructor *ClassMethod
	SourceMetadata
}

func (c *Class) GetType() NodeType                 { return ClassNode }
func (c *Class) GetSourceMetadata() SourceMetadata { return c.SourceMetadata }

type ClassMethod struct {
	Name     string
	Body     []Stmt
	Params   []string
	IsPublic bool
	SourceMetadata
}

func (c *ClassMethod) GetType() NodeType                 { return ClassMethodNode }
func (c *ClassMethod) GetSourceMetadata() SourceMetadata { return c.SourceMetadata }

type ClassProperty struct {
	Name     string
	Value    Expr
	IsPublic bool
	SourceMetadata
}

func (c *ClassProperty) GetType() NodeType                 { return ClassPropertyNode }
func (c *ClassProperty) GetSourceMetadata() SourceMetadata { return c.SourceMetadata }

type WhileLoop struct {
	Body      []Stmt
	Condition Expr
	SourceMetadata
}

func (w *WhileLoop) GetType() NodeType                 { return WhileLoopNode }
func (w *WhileLoop) GetSourceMetadata() SourceMetadata { return w.SourceMetadata }

type ReturnStmt struct {
	Value Expr
	SourceMetadata
}

func (r *ReturnStmt) GetType() NodeType                 { return ReturnStmtNode }
func (r *ReturnStmt) GetSourceMetadata() SourceMetadata { return r.SourceMetadata }

type BreakStmt struct {
	SourceMetadata
}

func (r *BreakStmt) GetType() NodeType                 { return BreakStmtNode }
func (r *BreakStmt) GetSourceMetadata() SourceMetadata { return r.SourceMetadata }

type ContinueStmt struct {
	SourceMetadata
}

func (r *ContinueStmt) GetType() NodeType                 { return ContinueStmtNode }
func (r *ContinueStmt) GetSourceMetadata() SourceMetadata { return r.SourceMetadata }

// Expressions

type VarAssignmentExpr struct {
	Assignee Expr
	Value    Expr
	SourceMetadata
}

func (v *VarAssignmentExpr) GetType() NodeType                 { return VarAssignmentExprNode }
func (v *VarAssignmentExpr) GetSourceMetadata() SourceMetadata { return v.SourceMetadata }

type BinaryExpr struct {
	LHS      Expr
	RHS      Expr
	Operator BinaryOperator
	SourceMetadata
}

func (b *BinaryExpr) GetType() NodeType                 { return BinaryExprNode }
func (b *BinaryExpr) GetSourceMetadata() SourceMetadata { return b.SourceMetadata }

type CompareExpr struct {
	LHS      Expr
	RHS      Expr
	Operator CompareOperator
	SourceMetadata
}

func (c *CompareExpr) GetType() NodeType                 { return CompareExprNode }
func (c *CompareExpr) GetSourceMetadata() SourceMetadata { return c.SourceMetadata }

type LogicalExpr struct {
	LHS      *Expr // Optional LHS for unary operators
	RHS      Expr
	Operator LogicalOperator
	SourceMetadata
}

func (l *LogicalExpr) GetType() NodeType                 { return LogicalExprNode }
func (l *LogicalExpr) GetSourceMetadata() SourceMetadata { return l.SourceMetadata }

type CallExpr struct {
	Args   []Expr
	Callee Expr
	SourceMetadata
}

func (c *CallExpr) GetType() NodeType                 { return CallExprNode }
func (c *CallExpr) GetSourceMetadata() SourceMetadata { return c.SourceMetadata }

type MemberExpr struct {
	Object   Expr
	Value    Expr
	Computed bool
	SourceMetadata
}

func (m *MemberExpr) GetType() NodeType                 { return MemberExprNode }
func (m *MemberExpr) GetSourceMetadata() SourceMetadata { return m.SourceMetadata }

type Identifier struct {
	Symbol string
	SourceMetadata
}

func (i *Identifier) GetType() NodeType                 { return IdentifierNode }
func (i *Identifier) GetSourceMetadata() SourceMetadata { return i.SourceMetadata }

type NumericLiteral struct {
	Value float64
	SourceMetadata
}

func (n *NumericLiteral) GetType() NodeType                 { return NumericLiteralNode }
func (n *NumericLiteral) GetSourceMetadata() SourceMetadata { return n.SourceMetadata }

type StringLiteral struct {
	Value string
	SourceMetadata
}

func (s *StringLiteral) GetType() NodeType                 { return StringLiteralNode }
func (s *StringLiteral) GetSourceMetadata() SourceMetadata { return s.SourceMetadata }

type Property struct {
	Key   string
	Value Expr
	SourceMetadata
}

func (p *Property) GetType() NodeType                 { return PropertyNode }
func (p *Property) GetSourceMetadata() SourceMetadata { return p.SourceMetadata }

type ObjectLiteral struct {
	Properties []Property
	SourceMetadata
}

func (o *ObjectLiteral) GetType() NodeType                 { return ObjectLiteralNode }
func (o *ObjectLiteral) GetSourceMetadata() SourceMetadata { return o.SourceMetadata }

type ArrayLiteral struct {
	Elements []Expr
	SourceMetadata
}

func (o *ArrayLiteral) GetType() NodeType                 { return ArrayLiteralNode }
func (o *ArrayLiteral) GetSourceMetadata() SourceMetadata { return o.SourceMetadata }
