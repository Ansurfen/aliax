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

type Ident struct {
	Name string
}

func (*Ident) exprNode() {}

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

type BlockStmt struct {
	List []Stmt
}

func (*BlockStmt) stmtNode() {}

type SwitchStmt struct {
	Cond  Expr
	Cases []*CaseStmt
}

func (*SwitchStmt) stmtNode() {}

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

type Comment struct {
	Text string
}

func (*Comment) stmtNode() {}
