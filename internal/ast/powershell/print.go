// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package psast

import (
	"aliax/internal/token/powershell"
	"fmt"
	"io"
	"strings"
)

func Print(node Node, w io.Writer) {
	print(w, node, "")
}

func print(w io.Writer, node Node, space string) {
	switch node := node.(type) {
	case *File:
		for _, s := range node.Stmts {
			print(w, s, space)
		}
	case *BlockStmt:
		fmt.Fprintf(w, " {\n")
		for _, s := range node.List {
			print(w, s, space+"  ")
		}
		fmt.Fprintln(w, space+"}")
	case *IfStmt:
		fmt.Fprintf(w, space+"if (%s)", exprString(node.Cond))
		print(w, node.Body, space)
		for node.Else != nil {
			switch el := node.Else.(type) {
			case *IfStmt:
				fmt.Fprintf(w, space+"elseif (%s)", exprString(el.Cond))
				print(w, el.Body, space)
				node.Else = el.Else
			case *BlockStmt:
				fmt.Fprint(w, space+"else")
				print(w, el, space)
				node.Else = nil
			}
		}
	case *ForStmt:
		fmt.Fprintf(w, space+"for (%s; %s; %s)", exprString(node.Init), exprString(node.Cond), exprString(node.Post))
		print(w, node.Body, space)
	case *SwitchStmt:
		if node.Regex {
			fmt.Fprintf(w, space+"switch -regex (%s) {\n", exprString(node.Cond))
		} else {
			fmt.Fprintf(w, space+"switch (%s) {\n", exprString(node.Cond))
		}
		for _, c := range node.Cases {
			fmt.Fprintf(w, space+"  %s", exprString(c.Cond))
			print(w, c.Body, space+"  ")
		}
		if node.Default != nil {
			fmt.Fprint(w, space+"  default")
			print(w, node.Default.Body, space+"  ")
		}
		fmt.Fprintln(w, space+"}")
	case *ExprStmt:
		fmt.Fprintln(w, space+exprString(node.X))
	case *AssignStmt:
		fmt.Fprintf(w, space+"%s = %s\n", exprString(node.Lhs), exprString(node.Rhs))
	case *CallStmt:
		recv := []string{}
		for _, r := range node.Recv {
			recv = append(recv, exprString(r))
		}
		if node.CallOp == token.AND {
			fmt.Fprintf(w, space+"& %s %s\n", exprString(node.Func), strings.Join(recv, " "))
		} else {
			fmt.Fprintf(w, space+"%s %s\n", exprString(node.Func), strings.Join(recv, " "))
		}
	case *Comment:
		fmt.Fprintf(w, space+"# %s\n", node.Text)
	}
}

func exprString(expr Expr) string {
	switch expr := expr.(type) {
	case *Ident:
		return expr.Name
	case *RefExpr:
		return fmt.Sprintf("$%s", exprString(expr.X))
	case *BasicExpr:
		switch expr.Kind {
		case token.STRING:
			return fmt.Sprintf(`"%s"`, expr.Value)
		default:
			return expr.Value
		}
	case *BinaryExpr:
		if expr.Op == token.DOUBLE_DOT {
			return fmt.Sprintf("%s..%s", exprString(expr.X), exprString(expr.Y))
		}
		return fmt.Sprintf("%s %s %s", exprString(expr.X), expr.Op, exprString(expr.Y))
	case *IndexExpr:
		return fmt.Sprintf("%s[%s]", exprString(expr.X), exprString(expr.Key))
	case *SelectorExpr:
		return fmt.Sprintf("%s.%s", exprString(expr.X), exprString(expr.Sel))
	case *IncDecExpr:
		return fmt.Sprintf("%s%s", exprString(expr.X), expr.Op)
	}
	return ""
}
