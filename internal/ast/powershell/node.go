// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package psast

import (
	"aliax/internal/token/powershell"
)

type Node interface{}

type Expr interface {
	Node
	exprNode()
}

type Comment struct {
	Text string
}

func (*Comment) stmtNode() {}

type CallStmt struct {
	CallOp token.Token
	Func   Expr
	Recv   []Expr
}

func (*CallStmt) stmtNode() {}

type BasicExpr struct {
	Kind  token.Token
	Value string
}

func (*BasicExpr) exprNode() {}

type IndexExpr struct {
	X   Expr
	Key Expr
}

func (*IndexExpr) exprNode() {}

type Ident struct {
	Name string
}

func (*Ident) exprNode() {}

type RefExpr struct {
	X Expr
}

func (*RefExpr) exprNode() {}

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

type Stmt interface {
	Node
	stmtNode()
}

type AssignStmt struct {
	Lhs Expr
	Rhs Expr
}

func (*AssignStmt) stmtNode() {}

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

type ForStmt struct {
	Init Expr
	Cond Expr
	Post Expr
	Body *BlockStmt
}

func (*ForStmt) stmtNode() {}

type SwitchStmt struct {
	Regex   bool
	Cond    Expr
	Cases   []*CaseStmt
	Default *CaseStmt
}

func (*SwitchStmt) stmtNode() {}

type CaseStmt struct {
	Cond Expr
	Body *BlockStmt
}

func (*CaseStmt) stmtNode() {}

type BlockStmt struct {
	List []Stmt
}

func (*BlockStmt) stmtNode() {}

type IncDecExpr struct {
	X  Expr
	Op token.Token
}

func (*IncDecExpr) exprNode() {}

var (
	Null  = &RefExpr{X: &Ident{Name: "null"}}
	True  = &RefExpr{X: &Ident{Name: "true"}}
	False = &RefExpr{X: &Ident{Name: "false"}}
)
