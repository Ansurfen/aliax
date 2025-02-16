// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package psast

import (
	"aliax/internal/token/powershell"
	"os"
	"testing"
)

func TestSwitch(t *testing.T) {
	Print(&SwitchStmt{
		Regex: true,
		Cond:  &RefExpr{X: &IndexExpr{X: &Ident{Name: "args"}, Key: &RefExpr{X: &Ident{Name: "i"}}}},
		Cases: []*CaseStmt{
			{Cond: &BasicExpr{Kind: token.STRING, Value: "dns|address"}, Body: &BlockStmt{
				List: []Stmt{
					&SwitchStmt{
						Regex: true,
						Cond:  &RefExpr{X: &IndexExpr{X: &Ident{Name: "args"}, Key: &RefExpr{X: &Ident{Name: "i"}}}},
						Cases: []*CaseStmt{
							{Cond: &BasicExpr{Kind: token.STRING, Value: "dns|address"}, Body: &BlockStmt{}},
						},
					}},
			}},
		},
		Default: &CaseStmt{
			Body: &BlockStmt{
				List: []Stmt{
					&AssignStmt{Lhs: &RefExpr{X: &Ident{"x"}}, Rhs: &RefExpr{X: &IndexExpr{X: &Ident{Name: "args"}, Key: &BinaryExpr{X: &RefExpr{X: &Ident{Name: "i"}}, Op: token.ADD, Y: &BasicExpr{Kind: token.NUMBER, Value: "1"}}}}},
					&ExprStmt{X: &IncDecExpr{X: &RefExpr{X: &Ident{Name: "i"}}, Op: token.Inc}},
				},
			},
		},
	}, os.Stdout)
}

func TestFor(t *testing.T) {
	Print(&ForStmt{
		Init: &BinaryExpr{X: &RefExpr{X: &Ident{Name: "i"}}, Op: token.ASSIGN, Y: &BasicExpr{Kind: token.NUMBER, Value: "0"}},
		Cond: &BinaryExpr{X: &RefExpr{X: &Ident{Name: "i"}}, Op: token.LT, Y: &RefExpr{X: &SelectorExpr{X: &Ident{Name: "args"}, Sel: &Ident{Name: "Length"}}}},
		Post: &IncDecExpr{X: &RefExpr{X: &Ident{Name: "i"}}, Op: token.Inc},
		Body: &BlockStmt{
			List: []Stmt{
				&ForStmt{
					Init: &BinaryExpr{X: &RefExpr{X: &Ident{Name: "i"}}, Op: token.ASSIGN, Y: &BasicExpr{Kind: token.NUMBER, Value: "0"}},
					Cond: &BinaryExpr{X: &RefExpr{X: &Ident{Name: "i"}}, Op: token.LT, Y: &RefExpr{X: &SelectorExpr{X: &Ident{Name: "args"}, Sel: &Ident{Name: "Length"}}}},
					Post: &IncDecExpr{X: &RefExpr{X: &Ident{Name: "i"}}, Op: token.Inc},
					Body: &BlockStmt{},
				},
			},
		},
	}, os.Stdout)
}

func TestIf(t *testing.T) {
	Print(&IfStmt{
		Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
		Body: &BlockStmt{
			List: []Stmt{
				&IfStmt{
					Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
					Body: &BlockStmt{},
					Else: &IfStmt{
						Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
						Body: &BlockStmt{},
						Else: &IfStmt{
							Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
							Body: &BlockStmt{},
							Else: &BlockStmt{
								List: []Stmt{
									&IfStmt{
										Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
										Body: &BlockStmt{},
										Else: &IfStmt{
											Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
											Body: &BlockStmt{},
											Else: &IfStmt{
												Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
												Body: &BlockStmt{},
												Else: &BlockStmt{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Else: &IfStmt{
			Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
			Body: &BlockStmt{
				List: []Stmt{
					&CallStmt{Func: &Ident{Name: "echo"}, Recv: []Expr{&BasicExpr{Kind: token.STRING, Value: "Hello World"}}},
					&CallStmt{Func: &Ident{Name: "exit"}},
				},
			},
			Else: &IfStmt{
				Cond: &BinaryExpr{X: Null, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
				Body: &BlockStmt{},
				Else: &BlockStmt{},
			},
		},
	}, os.Stdout)
}

func TestCall(t *testing.T) {
	Print(&CallStmt{
		CallOp: token.AND,
		Func:   &Ident{Name: "Write-Host"},
		Recv: []Expr{
			&Ident{Name: "$executable"},
			&Ident{Name: "-h"},
		},
	}, os.Stdout)
}
