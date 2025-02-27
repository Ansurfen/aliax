// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package psast

import (
	token "aliax/internal/token/powershell"
	"fmt"
)

// Node represents a generic node in the AST.
type Node interface{}

// Expr represents an expression node.
type Expr interface {
	Node
	exprNode()
	String() string
}

// Stmt represents a statement node.
type Stmt interface {
	Node
	stmtNode()
}

// Raw creates an expression from a raw script string.
func Raw(script string) Expr {
	return &Ident{Name: script}
}

// Comment represents a comment in the AST.
type Comment struct {
	Text string
}

func (*Comment) stmtNode() {}

// Docs creates a new comment node.
func Docs(text string) *Comment {
	return &Comment{Text: text}
}

type (
	// BasicExpr represents a basic expression with a token kind and value.
	BasicExpr struct {
		Kind  token.Token // The type of the basic expression (e.g., number, string, bool).
		Value string      // The value of the expression.
	}

	// IndexExpr represents an indexed expression (e.g., array or map access).
	IndexExpr struct {
		X   Expr // The target expression being indexed.
		Key Expr // The key used for indexing.
	}

	// Ident represents an identifier (variable or function name).
	Ident struct {
		Name string // The name of the identifier.
	}

	// RefExpr represents a reference to another expression.
	RefExpr struct {
		X Expr // The expression being referenced.
	}

	// BinaryExpr represents a binary operation expression.
	BinaryExpr struct {
		X  Expr        // The left-hand side expression.
		Op token.Token // The binary operator.
		Y  Expr        // The right-hand side expression.
	}

	// SelectorExpr represents a field or method selector (e.g., obj.field).
	SelectorExpr struct {
		X   Expr // The expression being accessed.
		Sel Expr // The selected field or method.
	}

	// CaseStmt represents a case statement in a switch.
	CaseStmt struct {
		Cond Expr       // The case condition.
		Body *BlockStmt // The case body.
	}

	// IncDecExpr represents an increment or decrement expression.
	IncDecExpr struct {
		X  Expr        // The expression being incremented or decremented.
		Op token.Token // The operator (increment or decrement).
	}
)

func (*BasicExpr) exprNode()    {}
func (*IndexExpr) exprNode()    {}
func (*Ident) exprNode()        {}
func (*RefExpr) exprNode()      {}
func (*BinaryExpr) exprNode()   {}
func (*SelectorExpr) exprNode() {}
func (*CaseStmt) stmtNode()     {}
func (*IncDecExpr) exprNode()   {}

func (e *BasicExpr) String() string {
	switch e.Kind {
	case token.STRING:
		return fmt.Sprintf(`"%s"`, e.Value)
	default:
		return e.Value
	}
}

func (e *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", e.X, e.Key)
}

func (e *Ident) String() string {
	return e.Name
}

func (e *RefExpr) String() string {
	return fmt.Sprintf("$%s", e.X)
}

func (e *BinaryExpr) String() string {
	if e.Op == token.DOUBLE_DOT {
		return fmt.Sprintf("%s..%s", e.X, e.Y)
	}
	return fmt.Sprintf("%s %s %s", e.X, e.Op, e.Y)
}

func (e *SelectorExpr) String() string {
	return fmt.Sprintf("%s.%s", e.X, e.Sel)
}

func (e *IncDecExpr) String() string {
	return fmt.Sprintf("%s%s", e.X, e.Op)
}

// Number creates a number expression.
func Number(n int) *BasicExpr {
	return &BasicExpr{
		Kind:  token.NUMBER,
		Value: fmt.Sprintf("%d", n),
	}
}

// Bool creates a boolean expression.
func Bool(b bool) *BasicExpr {
	return &BasicExpr{
		Kind:  token.BOOL,
		Value: fmt.Sprintf("%t", b),
	}
}

// String creates a string expression.
func String(s string) *BasicExpr {
	return &BasicExpr{
		Kind:  token.STRING,
		Value: s,
	}
}

// IndexExpression creates a new index expression.
func IndexExpression(x, key Expr) *IndexExpr {
	return &IndexExpr{
		X:   x,
		Key: key,
	}
}

// Identifier creates a new identifier.
func Identifier(name string) *Ident {
	return &Ident{Name: name}
}

var (
	NULL  = &RefExpr{X: &Ident{Name: "null"}}
	TRUE  = &RefExpr{X: &Ident{Name: "true"}}
	FALSE = &RefExpr{X: &Ident{Name: "false"}}
)

// RefRaw creates a reference expression from a raw string.
func RefRaw(name string) *RefExpr {
	return &RefExpr{X: &Ident{Name: name}}
}

// RefExpression creates a reference expression.
func RefExpression(x Expr) *RefExpr {
	return &RefExpr{X: x}
}

// BinaryExpression creates a new binary operation expression.
func BinaryExpression(x Expr, op token.Token, y Expr) *BinaryExpr {
	return &BinaryExpr{
		X:  x,
		Op: op,
		Y:  y,
	}
}

// SelectorExpression creates a new selector expression.
func SelectorExpression(x, sel Expr) *SelectorExpr {
	return &SelectorExpr{
		X:   x,
		Sel: sel,
	}
}

// CaseStatement creates a new case statement.
func CaseStatement(cond Expr) *CaseStmt {
	return &CaseStmt{
		Cond: cond,
		Body: &BlockStmt{},
	}
}

// IncDecExpression creates an increment or decrement expression.
func IncDecExpression(name string, inc bool) *IncDecExpr {
	if inc {
		return &IncDecExpr{X: RefRaw(name), Op: token.Inc}
	}
	return &IncDecExpr{X: RefRaw(name), Op: token.Dec}
}

type (
	// CallStmt represents a function call statement.
	CallStmt struct {
		Op   token.Token // The token representing the call operator or the dot sourcing.
		Func Expr        // The function being called.
		Recv []Expr      // The arguments passed to the function.
	}

	// AssignStmt represents an assignment statement.
	AssignStmt struct {
		Lhs Expr // The left-hand side expression.
		Rhs Expr // The right-hand side expression.
	}

	// ExprStmt represents an expression statement.
	ExprStmt struct {
		X Expr // The expression being evaluated.
	}

	// IfStmt represents an if statement.
	IfStmt struct {
		Cond Expr       // The condition expression.
		Body *BlockStmt // The body of the if statement.
		Else Stmt       // The else clause (optional).
	}

	// ForStmt represents a for-loop statement.
	ForStmt struct {
		Init Expr       // Initialization expression (optional).
		Cond Expr       // Condition expression.
		Post Expr       // Post-expression (executed after each loop iteration).
		Body *BlockStmt // The loop body.
	}

	// SwitchStmt represents a switch statement.
	SwitchStmt struct {
		Mode    MatchMode
		Cond    Expr        // The condition expression.
		Cases   []*CaseStmt // List of case statements.
		Default *CaseStmt   // Default case (optional).
	}

	// BlockStmt represents a block of statements.
	BlockStmt struct {
		List []Stmt // List of statements in the block.
	}
)

func (*CallStmt) stmtNode()   {}
func (*ExprStmt) stmtNode()   {}
func (*AssignStmt) stmtNode() {}
func (*IfStmt) stmtNode()     {}
func (*ForStmt) stmtNode()    {}
func (*SwitchStmt) stmtNode() {}
func (*BlockStmt) stmtNode()  {}

// AssignStatement creates an assignment statement.
func AssignStatement(lhs, rhs Expr) *AssignStmt {
	return &AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
	}
}

// IfStatement creates an if statement with an empty body.
func IfStatement() *IfStmt {
	return &IfStmt{Body: &BlockStmt{}}
}

// ForStatement creates a for-loop statement.
func ForStatement(init, cond, post Expr) *ForStmt {
	return &ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: &BlockStmt{},
	}
}

type MatchMode uint8

const (
	MatchModeNone MatchMode = iota
	MatchModeWildcard
	MatchModeExact
	MatchModeCaseSensitive
	MatchModeFile
	MatchModeParallel
	MatchModeRegex
)

func (m MatchMode) String() string {
	return []string{"", "-Wildcard", "-Exact", "-CaseSensitive", "-File", "-Parallel", "-Regex"}[m]
}

// SwtichStatement creates a new switch statement.
func SwtichStatement(mode MatchMode, cond Expr, cases []*CaseStmt, default_ *CaseStmt) *SwitchStmt {
	return &SwitchStmt{
		Mode:    mode,
		Cond:    cond,
		Cases:   cases,
		Default: default_,
	}
}

// CallStatement creates a new function call statement.
func CallStatement(op token.Token, fn string, recv ...Expr) *CallStmt {
	return &CallStmt{
		Op:   op,
		Func: Identifier(fn),
		Recv: recv,
	}
}

func BlockStatement(stmts ...Stmt) *BlockStmt {
	return &BlockStmt{List: stmts}
}

// Append adds statements to a block.
func (b *BlockStmt) Append(stmts ...Stmt) {
	b.List = append(b.List, stmts...)
}

// File represents an AST root containing a list of statements.
type File struct {
	Stmts []Stmt // List of statements in the file.
}

func (f *File) Append(stmts ...Stmt) {
	f.Stmts = append(f.Stmts, stmts...)
}
