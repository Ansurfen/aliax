// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package bashast

import (
	token "aliax/internal/token/bash"
	"fmt"
)

// Node represents any node in the abstract syntax tree (AST).
type Node interface{}

// Expr is the interface for all expression nodes in the AST.
// Any type that implements `exprNode()` is considered an expression node.
type Expr interface {
	Node
	exprNode()
}

// Raw creates a raw expression from a string (like an identifier).
func Raw(script string) Expr {
	return &Ident{Name: script}
}

// RawStmt creates a raw statement from a string (like an expression statement).
func RawStmt(s string) Stmt {
	return &ExprStmt{X: &Ident{Name: s}}
}

// BinaryExpr represents a binary expression with a left operand (X), operator (Op), and right operand (Y).
type BinaryExpr struct {
	X  Expr
	Op token.Token
	Y  Expr
}

func (*BinaryExpr) exprNode() {}

// BinaryExpression creates a new binary expression with the given operands and operator.
func BinaryExpression(x Expr, op token.Token, y Expr) *BinaryExpr {
	return &BinaryExpr{X: x, Op: op, Y: y}
}

// SelectorExpr represents an expression where one expression accesses a field or method of another.
type SelectorExpr struct {
	X   Expr
	Sel Expr
}

func (*SelectorExpr) exprNode() {}

// File represents a collection of statements (like a program or a script).
type File struct {
	Stmts []Stmt
}

// RefExpr represents a reference to an expression (e.g., variable reference).
type RefExpr struct {
	X Expr
}

func (*RefExpr) exprNode() {}

func RefExpression(x Expr) *RefExpr {
	return &RefExpr{X: x}
}

func RefRaw(name string) *RefExpr {
	return &RefExpr{X: &Ident{Name: name}}
}

// IncDecExpr represents an increment or decrement operation on an expression.
type IncDecExpr struct {
	X  Expr
	Op token.Token
}

func (*IncDecExpr) exprNode() {}

// IncDecExpression creates an increment or decrement expression for the given operand.
func IncDecExpression(x Expr, inc bool) *IncDecExpr {
	if inc {
		return &IncDecExpr{X: x, Op: token.Inc}
	}
	return &IncDecExpr{X: x, Op: token.Dec}
}

// IndexExpr represents an expression where an index is used to access an element (e.g., array or slice indexing).
type IndexExpr struct {
	X   Expr
	Key Expr
}

func (*IndexExpr) exprNode() {}

// BasicExpr represents a basic expression with a type and a value.
type BasicExpr struct {
	Kind  token.Token
	Value string
}

func (*BasicExpr) exprNode() {}

// Number creates a new basic expression representing a number.
func Number(n int) *BasicExpr {
	return &BasicExpr{
		Kind:  token.NUMBER,
		Value: fmt.Sprintf("%d", n),
	}
}

// Bool creates a new basic expression representing a boolean value.
func Bool(b bool) *BasicExpr {
	return &BasicExpr{
		Kind:  token.BOOL,
		Value: fmt.Sprintf("%t", b),
	}
}

// String creates a new basic expression representing a string.
func String(s string) *BasicExpr {
	return &BasicExpr{
		Kind:  token.STRING,
		Value: s,
	}
}

var (
	TRUE  = &BasicExpr{Kind: token.BOOL, Value: "true"}
	FALSE = &BasicExpr{Kind: token.BOOL, Value: "false"}
)

// Ident represents an identifier (like a variable or function name).
type Ident struct {
	Name string
}

func (*Ident) exprNode() {}

// Identifier creates a new identifier expression with the given name.
func Identifier(name string) *Ident {
	return &Ident{Name: name}
}

// Stmt is the interface for all statement nodes in the AST.
// Any type that implements `stmtNode()` is considered a statement node.
type Stmt interface {
	Node
	stmtNode()
}

// IfStmt represents an `if` statement, which has a condition, a body, and an optional else branch.
type IfStmt struct {
	Cond Expr
	Body *BlockStmt
	Else Stmt
}

func (*IfStmt) stmtNode() {}

// IfStatement creates a new if statement with an empty body.
func IfStatement() *IfStmt {
	return &IfStmt{Body: &BlockStmt{}}
}

// ForStmt represents a `for` loop statement, with initialization, condition, post, and body.
type ForStmt struct {
	Init Expr
	Cond Expr
	Post Expr
	Body *BlockStmt
}

func (*ForStmt) stmtNode() {}

// ForStatement creates a new `for` statement with the given initialization, condition, and post-expressions.
func ForStatement(init, cond, post Expr) *ForStmt {
	return &ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: &BlockStmt{},
	}
}

// ExprStmt represents a statement that contains a single expression.
type ExprStmt struct {
	X Expr
}

func (*ExprStmt) stmtNode() {}

// BlockStmt represents a block of statements enclosed in braces `{ ... }`.
type BlockStmt struct {
	List []Stmt
}

func (*BlockStmt) stmtNode() {}

// BlockStatement creates a new block statement with the given list of statements.
func BlockStatement(stmts ...Stmt) *BlockStmt {
	return &BlockStmt{List: stmts}
}

// Append adds more statements to the end of the block.
func (b *BlockStmt) Append(stmts ...Stmt) {
	b.List = append(b.List, stmts...)
}

// SwitchStmt represents a `switch` statement with cases and a default block.
type SwitchStmt struct {
	Cond    Expr
	Cases   []*CaseStmt
	Default *CaseStmt
}

func (*SwitchStmt) stmtNode() {}

// SetDefault sets the default case for the switch statement.
func (s *SwitchStmt) SetDefault(b *BlockStmt) {
	s.Default = &CaseStmt{Body: b}
}

// CaseStmt represents a `case` statement in a switch, with a condition and a body.
type CaseStmt struct {
	Cond Expr
	Body *BlockStmt
}

func (*CaseStmt) stmtNode() {}

// CaseStatement creates a new case statement with the given condition.
func CaseStatement(cond Expr) *CaseStmt {
	return &CaseStmt{
		Cond: cond,
		Body: &BlockStmt{},
	}
}

// AssignStmt represents an assignment statement with a left-hand side (Lhs) and a right-hand side (Rhs).
type AssignStmt struct {
	Lhs Expr
	Rhs Expr
}

func (*AssignStmt) stmtNode() {}

// AssignStatement creates a new assignment statement with the given left and right expressions.
func AssignStatement(lhs, rhs Expr) *AssignStmt {
	return &AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
	}
}

// CallStmt represents a function call statement with a function name and arguments.
type CallStmt struct {
	Func Expr
	Recv []Expr
}

func (*CallStmt) stmtNode() {}

// CallStatement creates a new function call statement with the given function name and arguments.
func CallStatement(name string, args ...string) *CallStmt {
	recv := []Expr{}
	for _, a := range args {
		recv = append(recv, &Ident{Name: a})
	}
	return &CallStmt{
		Func: &Ident{Name: name},
		Recv: recv,
	}
}

// Comment represents a comment in the code.
type Comment struct {
	Text string
}

func (*Comment) stmtNode() {}

// Docs creates a new comment node with the given text.
func Docs(text string) *Comment {
	return &Comment{Text: text}
}
