// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package psast

import (
	token "aliax/internal/token/powershell"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier(t *testing.T) {
	id := Identifier("myVar")
	assert.Equal(t, "myVar", id.Name)
}

func TestNumber(t *testing.T) {
	num := Number(42)
	assert.Equal(t, token.NUMBER, num.Kind)
	assert.Equal(t, "42", num.Value)
}

func TestBool(t *testing.T) {
	b := Bool(true)
	assert.Equal(t, token.BOOL, b.Kind)
	assert.Equal(t, "true", b.Value)
}

func TestString(t *testing.T) {
	str := String("hello")
	assert.Equal(t, token.STRING, str.Kind)
	assert.Equal(t, "hello", str.Value)
}

func TestAssignStatement(t *testing.T) {
	lhs := Identifier("x")
	rhs := Number(10)
	assign := AssignStatement(lhs, rhs)

	assert.Equal(t, lhs, assign.Lhs)
	assert.Equal(t, rhs, assign.Rhs)
}

func TestBinaryExpression(t *testing.T) {
	left := Number(2)
	right := Number(3)
	expr := BinaryExpression(left, token.ADD, right)

	assert.Equal(t, left, expr.X)
	assert.Equal(t, right, expr.Y)
	assert.Equal(t, token.ADD, expr.Op)
}

func TestPrint(t *testing.T) {
	var buf strings.Builder

	Print(Identifier("foo"), &buf)
	assert.Equal(t, "foo", buf.String())
	buf.Reset()

	Print(Number(42), &buf)
	assert.Equal(t, "42", buf.String())
	buf.Reset()

	assign := AssignStatement(Identifier("x"), Number(100))
	Print(assign, &buf)
	assert.Equal(t, "x = 100", buf.String())
	buf.Reset()

	binary := BinaryExpression(Number(2), token.ADD, Number(3))
	Print(binary, &buf)
	assert.Equal(t, "(2 + 3)", buf.String())
}

func TestSwitch(t *testing.T) {
	Print(&SwitchStmt{
		Mode: MatchModeRegex,
		Cond: &RefExpr{X: &IndexExpr{X: &Ident{Name: "args"}, Key: &RefExpr{X: &Ident{Name: "i"}}}},
		Cases: []*CaseStmt{
			{Cond: &BasicExpr{Kind: token.STRING, Value: "dns|address"}, Body: &BlockStmt{
				List: []Stmt{
					&SwitchStmt{
						Mode: MatchModeRegex,
						Cond: &RefExpr{X: &IndexExpr{X: &Ident{Name: "args"}, Key: &RefExpr{X: &Ident{Name: "i"}}}},
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
		Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
		Body: &BlockStmt{
			List: []Stmt{
				&IfStmt{
					Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
					Body: &BlockStmt{},
					Else: &IfStmt{
						Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
						Body: &BlockStmt{},
						Else: &IfStmt{
							Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
							Body: &BlockStmt{},
							Else: &BlockStmt{
								List: []Stmt{
									&IfStmt{
										Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
										Body: &BlockStmt{},
										Else: &IfStmt{
											Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
											Body: &BlockStmt{},
											Else: &IfStmt{
												Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
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
			Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
			Body: &BlockStmt{
				List: []Stmt{
					&CallStmt{Func: &Ident{Name: "echo"}, Recv: []Expr{&BasicExpr{Kind: token.STRING, Value: "Hello World"}}},
					&CallStmt{Func: &Ident{Name: "exit"}},
				},
			},
			Else: &IfStmt{
				Cond: &BinaryExpr{X: NULL, Op: token.NE, Y: &RefExpr{X: &Ident{Name: "i"}}},
				Body: &BlockStmt{},
				Else: &BlockStmt{},
			},
		},
	}, os.Stdout)
}

func TestCall(t *testing.T) {
	Print(&CallStmt{
		Op:   token.BITAND,
		Func: &Ident{Name: "Write-Host"},
		Recv: []Expr{
			&Ident{Name: "$executable"},
			&Ident{Name: "-h"},
		},
	}, os.Stdout)
}
