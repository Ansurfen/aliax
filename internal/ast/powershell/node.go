// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package psast

import (
	token "aliax/internal/token/powershell"
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

type Comment struct {
	Text string
}

func (*Comment) stmtNode() {}

func Docs(text string) *Comment {
	return &Comment{Text: text}
}

type CallStmt struct {
	CallOp token.Token
	Func   Expr
	Recv   []Expr
}

func (*CallStmt) stmtNode() {}

func CallStatement(op token.Token, fn string, recv ...Expr) *CallStmt {
	return &CallStmt{
		CallOp: op,
		Func:   Identifier(fn),
		Recv:   recv,
	}
}

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

type IndexExpr struct {
	X   Expr
	Key Expr
}

func (*IndexExpr) exprNode() {}

func IndexExpression(x, key Expr) *IndexExpr {
	return &IndexExpr{
		X:   x,
		Key: key,
	}
}

type Ident struct {
	Name string
}

func (*Ident) exprNode() {}

func Identifier(name string) *Ident {
	return &Ident{Name: name}
}

type RefExpr struct {
	X Expr
}

func (*RefExpr) exprNode() {}

func RefRaw(name string) *RefExpr {
	return &RefExpr{X: &Ident{Name: name}}
}

func RefExpression(x Expr) *RefExpr {
	return &RefExpr{X: x}
}

type BinaryExpr struct {
	X  Expr
	Op token.Token
	Y  Expr
}

func (*BinaryExpr) exprNode() {}

func BinaryExpression(x Expr, op token.Token, y Expr) *BinaryExpr {
	return &BinaryExpr{
		X:  x,
		Op: op,
		Y:  y,
	}
}

type SelectorExpr struct {
	X   Expr
	Sel Expr
}

func (*SelectorExpr) exprNode() {}

func SelectorExpression(x, sel Expr) *SelectorExpr {
	return &SelectorExpr{
		X:   x,
		Sel: sel,
	}
}

type File struct {
	Stmts []Stmt
}

type Stmt interface {
	Node
	stmtNode()
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

type ExprStmt struct {
	X Expr
}

func (*ExprStmt) stmtNode() {}

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

type SwitchStmt struct {
	Regex   bool
	Cond    Expr
	Cases   []*CaseStmt
	Default *CaseStmt
}

func (*SwitchStmt) stmtNode() {}

func SwtichStatement(reg bool, cond Expr, cases []*CaseStmt, default_ *CaseStmt) *SwitchStmt {
	return &SwitchStmt{
		Regex:   reg,
		Cond:    cond,
		Cases:   cases,
		Default: default_,
	}
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

type BlockStmt struct {
	List []Stmt
}

func (*BlockStmt) stmtNode() {}

func (b *BlockStmt) Append(stmts ...Stmt) {
	b.List = append(b.List, stmts...)
}

type IncDecExpr struct {
	X  Expr
	Op token.Token
}

func (*IncDecExpr) exprNode() {}

func IncDecExpression(name string, inc bool) *IncDecExpr {
	if inc {
		return &IncDecExpr{X: RefRaw(name), Op: token.Inc}
	}
	return &IncDecExpr{X: RefRaw(name), Op: token.Dec}
}

var (
	Null  = &RefExpr{X: &Ident{Name: "null"}}
	True  = &RefExpr{X: &Ident{Name: "true"}}
	False = &RefExpr{X: &Ident{Name: "false"}}
)
