// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package bashast

import "aliax/internal/token/bash"

type Node interface{}

type Expr interface {
	Node
	exprNode()
}

type BinaryExpr struct {
	X  Expr
	Op token.Token
	Y  Expr
}

func (*BinaryExpr) exprNode() {}

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

var (
	TRUE  = &BasicExpr{Kind: token.BOOL, Value: "true"}
	FALSE = &BasicExpr{Kind: token.BOOL, Value: "false"}
)

type Ident struct {
	Name string
}

func (*Ident) exprNode() {}

func NewIdent(name string) *Ident {
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

func NewIfStmt() *IfStmt {
	return &IfStmt{Body: &BlockStmt{}}
}

type ForStmt struct {
	Init Expr
	Cond Expr
	Post Expr
	Body *BlockStmt
}

func (*ForStmt) stmtNode() {}

type ExprStmt struct {
	X Expr
}

func (*ExprStmt) stmtNode() {}

func NewExprStmt(s string) *ExprStmt {
	return &ExprStmt{X: &Ident{Name: s}}
}

type BlockStmt struct {
	List []Stmt
}

func (*BlockStmt) stmtNode() {}

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

type AssignStmt struct {
	Lhs Expr
	Rhs Expr
}

func (*AssignStmt) stmtNode() {}

type CallStmt struct {
	Func Expr
	Recv []Expr
}

func (*CallStmt) stmtNode() {}

func NewCallStmt(name string, args ...string) *CallStmt {
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
