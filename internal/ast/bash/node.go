// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package bashast

import (
	token "aliax/internal/token/bash"
	"fmt"
)

type Node interface{}

type Expr interface {
	Node
	exprNode()
}

func Raw(script string) Expr {
	return &Ident{Name: script}
}

func RawStmt(s string) Stmt {
	return &ExprStmt{X: &Ident{Name: s}}
}

type BinaryExpr struct {
	X  Expr
	Op token.Token
	Y  Expr
}

func (*BinaryExpr) exprNode() {}

func BinaryExpression(x Expr, op token.Token, y Expr) *BinaryExpr {
	return &BinaryExpr{X: x, Op: op, Y: y}
}

type SelectorExpr struct {
	X   Expr
	Sel Expr
}

func (*SelectorExpr) exprNode() {}

type File struct {
	Stmts []Stmt
}

type RefExpr struct {
	X Expr
}

func (*RefExpr) exprNode() {}

type IncDecExpr struct {
	X  Expr
	Op token.Token
}

func (*IncDecExpr) exprNode() {}

func IncDecExpression(x Expr, inc bool) *IncDecExpr {
	if inc {
		return &IncDecExpr{X: x, Op: token.Inc}
	}
	return &IncDecExpr{X: x, Op: token.Dec}
}

type IndexExpr struct {
	X   Expr
	Key Expr
}

func (*IndexExpr) exprNode() {}

type BasicExpr struct {
	Kind  token.Token
	Value string
}

func (*BasicExpr) exprNode() {}

func Number(n int) *BasicExpr {
	return &BasicExpr{
		Kind:  token.NUMBER,
		Value: fmt.Sprintf("%d", n),
	}
}

func Bool(b bool) *BasicExpr {
	return &BasicExpr{
		Kind:  token.BOOL,
		Value: fmt.Sprintf("%t", b),
	}
}

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

type Ident struct {
	Name string
}

func (*Ident) exprNode() {}

func Identifier(name string) *Ident {
	return &Ident{Name: name}
}

type Stmt interface {
	Node
	stmtNode()
}

type IfStmt struct {
	Cond Expr
	Body *BlockStmt
	Else Stmt
}

func (*IfStmt) stmtNode() {}

func IfStatement() *IfStmt {
	return &IfStmt{Body: &BlockStmt{}}
}

type ForStmt struct {
	Init Expr
	Cond Expr
	Post Expr
	Body *BlockStmt
}

func (*ForStmt) stmtNode() {}

func ForStatement(init, cond, post Expr) *ForStmt {
	return &ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: &BlockStmt{},
	}
}

type ExprStmt struct {
	X Expr
}

func (*ExprStmt) stmtNode() {}

type BlockStmt struct {
	List []Stmt
}

func (*BlockStmt) stmtNode() {}

func BlockStatement(stmts ...Stmt) *BlockStmt {
	return &BlockStmt{List: stmts}
}

func (b *BlockStmt) Append(stmts ...Stmt) {
	b.List = append(b.List, stmts...)
}

type SwitchStmt struct {
	Cond    Expr
	Cases   []*CaseStmt
	Default *CaseStmt
}

func (*SwitchStmt) stmtNode() {}

func (s *SwitchStmt) SetDefault(b *BlockStmt) {
	s.Default = &CaseStmt{Body: b}
}

type CaseStmt struct {
	Cond Expr
	Body *BlockStmt
}

func (*CaseStmt) stmtNode() {}

func CaseStatement(cond Expr) *CaseStmt {
	return &CaseStmt{
		Cond: cond,
		Body: &BlockStmt{},
	}
}

type AssignStmt struct {
	Lhs Expr
	Rhs Expr
}

func (*AssignStmt) stmtNode() {}

func AssignStatement(lhs, rhs Expr) *AssignStmt {
	return &AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
	}
}

type CallStmt struct {
	Func Expr
	Recv []Expr
}

func (*CallStmt) stmtNode() {}

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

type Comment struct {
	Text string
}

func (*Comment) stmtNode() {}
